# Server
```
cd system
go test
```

# Client
```
$ curl http://127.0.0.1:8010/demo/DescribeDemo
before demo action
DescribeDemo
after demo action

$ curl http://127.0.0.1:8010/demo/Redirect
<a href="https://baidu.com">Found</a>
before demo action
```

# 思路
1. 定义 Controller。每个 Controller 中有多个 Action，每个 Action 的第 1 个参数为 Controller 自身，第 2 个参数为 ActionContext。ActionContext 中主要记录 Request 和 Response。
2. 注册路由。router.RegisRoutes 自动获取 Controller 下的所有 Action 并注册到路由表中。
3. 查找路由。router.FindRoute 从 URL 中解析出 Controller 名称和 Action 名称，根据 Controller 名称查找路由表，在找到的路由项中继续根据 Action 名称查找 Action。
4. 请求处理。system.ServeHTTP 用查找到的 Action 对请求进行处理。

# 参考
https://www.jianshu.com/p/25015167e21c
