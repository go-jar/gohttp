# 1 Example

## 1.1 Server
```
$ cd gohttp
$ go mod init gohttp
$ go build -o server server_demo.go
$ ./server

Hi, Client! Your data is: 123
destruct demo context
Hi, Client! Your data is: 123
destruct demo context
destruct demo context
a=1&b=2&c=3
destruct demo context

2020/10/09 09:43:56 [pid: 32339] Receive signal SIGUSR2. To restart the server gracefully.
2020/10/09 09:43:56 [pid: 32339] A new process [32759] has been started successfully. To shut down the server gracefully.
2020/10/09 09:43:56 [pid: 32339] The http server has been shut down successfully.
2020/10/09 09:43:56 [pid: 32339] Waiting for connections to be closed.
2020/10/09 09:43:56 [pid: 32339] All connections has been closed.

2020/10/09 09:45:48 [pid: 32759] Receive signal SIGUSR2. To restart the server gracefully.
2020/10/09 09:45:48 [pid: 32759] A new process [832] has been started successfully. To shut down the server gracefully.
2020/10/09 09:45:48 [pid: 32759] The http server has been shut down successfully.
2020/10/09 09:45:48 [pid: 32759] Waiting for connections to be closed.
2020/10/09 09:45:48 [pid: 32759] All connections has been closed.
```

## 1.2 Client
```
$ curl http://127.0.0.1:8010/demo/DescribeDemo
before demo action
DescribeDemo
after demo action

$ curl http://127.0.0.1:8010/demo/Redirect
<a href="https://baidu.com">Found</a>
before demo action
```

```
$ cd gohttp/httpclient
$ go test

[Debug] [2020-12-13 22:31:03]   -       [HttpClient]    Host: 127.0.0.1:8010    URL: http://127.0.0.1:8010/demo/args/123        TimeDuration: 986.798µs StatusCode: 200
Hi, Client! Your data is: 123
before demo action
after demo action
 986.798µs
[Debug] [2020-12-13 22:31:03]   -       [HttpClient]    Host: 127.0.0.1:8010    URL: http://127.0.0.1:8010/test/args/123        TimeDuration: 566.368µs StatusCode: 200
Hi, Client! Your data is: 123
before demo action
after demo action
 566.368µs
[Debug] [2020-12-13 22:31:03]   -       [HttpClient]    Host: 127.0.0.1:8010    URL: http://127.0.0.1:8010/demo/DescribeDemo    TimeDuration: 606.191µs StatusCode: 200
before demo action
DescribeDemo
after demo action
 606.191µs
[Debug] [2020-12-13 22:31:03]   -       [HttpClient]    Host: 127.0.0.1:8010    URL: http://127.0.0.1:8010/demo/ProcessPost     TimeDuration: 482.174µs StatusCode: 200
Hi, Client! Your data is: a=1&b=2&c=3
before demo action
after demo action
 482.174µs
PASS
ok      gohttp/httpclient       0.006s

$ ps -ef | grep example
root     32339 29709  0 09:42 pts/1    00:00:00 ./example
root     32621 30048  0 09:43 pts/6    00:00:00 grep --color=auto example

$ kill -USR2 32339

$ ps -ef | grep example
root       650 30048  0 09:45 pts/6    00:00:00 grep --color=auto example
root     32759     1  0 09:43 pts/1    00:00:00 ./example

$ kill -USR2 32759
```

# 2 gracehttp

## 2.1 热重启
进程在不关闭监听端口的情况下重启，重启期间所有的请求能被正确处理。

## 2.2 步骤
1. 父进程监听信号。
    - 如果是 SIGTERM，则关闭；
    - 如果是 SIGUSR2，则优雅重启，转步骤 2。
2. 父进程收到 SIGUSR2 信号后 fork 子进程，将服务监听的 socket 文件描述符传递给子进程。
    - 可以采用 syscall.ForkExec，使用完全相同的参数启动子进程。exec 是在调用进程内部执行一个可执行文件。
    - 将优雅启动标识也传递给子进程，因为需要告诉子进程这是优雅重启，子进程应该复用当前 socket，而不是打开一个新的 socket。
3. 子进程监听父进程的 socket，此时父进程和子进程都可以接收请求。
4. 子进程启动成功后，父进程优雅退出，即停止接收新的连接，等待旧连接处理完成（或超时）。
    - 父进程可通过 sync.WaitGroup 跟踪所有打开的连接数。每当新连接到来时（Accept），增加计数；连接关闭时（close），减小计数。
    - 但是 Go 中的 http.Server.Shutdown 已经实现了优雅退出。
5. 父进程退出，重启完成。

# 3 router

1. 定义 Controller。每个 Controller 中有多个 Action，每个 Action 的第 1 个参数为 Controller 自身，第 2 个参数为 ActionContext。ActionContext 中主要记录 Request 和 Response。
2. 注册路由。
    - 一般路由注册：router.RegisterRoutes 自动获取 Controller 下的所有 Action 并注册到路由表中。
    - 注册指定路由：注册 Controller 中的 Action 时，将需要匹配的 path 规则也注册好，该规则主要是查找 url 中除了 Controller 和 Action 以外的参数。
3. 查找路由。
    - 一般路由查找：router.FindRoute 从 URL 中解析出 Controller 名称和 Action 名称，根据 Controller 名称查找路由表，在找到的路由项中继续根据 Action 名称查找 Action。
    - 查找指定路由：查找匹配指定 path 规则的路由。
4. 请求处理。system.ServeHTTP 用查找到的 Action 对请求进行处理。

# 参考
- https://www.jianshu.com/p/25015167e21c
- https://github.com/goinbox/gohttp
