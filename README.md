# 1 Example

## 1.1 Server
```
$ cd system
$ go test

destruct demo context
destruct demo context
2020/10/08 15:36:58 [pid: 31061] Receive signal SIGUSR2. To restart the server gracefully.
2020/10/08 15:36:58 [pid: 31061] A new process [31129] has been started successfully. To shut down the server gracefully.
2020/10/08 15:36:58 [pid: 31061] The http server has been shut down successfully.
2020/10/08 15:36:58 [pid: 31061] Waiting for connections to be closed.
2020/10/08 15:36:58 [pid: 31061] All connections has been closed.
PASS
ok  	gohttp/system	19.913s
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

$ ps -ef | grep go
root     31004 23568  2 15:36 pts/1    00:00:00 go test
root     31061 31004  0 15:36 pts/1    00:00:00 /tmp/go-build404490548/b001/system.test -test.timeout=10m0s
root     31094 23887  0 15:36 pts/6    00:00:00 grep --color=auto go

$ kill -USR2 31061
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
2. 注册路由。router.RegisRoutes 自动获取 Controller 下的所有 Action 并注册到路由表中。
3. 查找路由。router.FindRoute 从 URL 中解析出 Controller 名称和 Action 名称，根据 Controller 名称查找路由表，在找到的路由项中继续根据 Action 名称查找 Action。
4. 请求处理。system.ServeHTTP 用查找到的 Action 对请求进行处理。

# 参考
- https://www.jianshu.com/p/25015167e21c
- https://github.com/goinbox/gohttp
