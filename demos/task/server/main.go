package main

import (
	jobs "example/demos/task"
	. "my-zinx/logging"
	"my-zinx/server"
)

/*
背景：以往的HTTP服务器，如果要实现一个简单的加减乘除计算器，可以这么实现：
1. 客户端发起GET请求，携带请求参数，例如：/add?a=1&b=2
2. 服务器接收到请求，路由add后解析参数，进行计算，返回结果。
3. 客户端收到结果，断开连接。
因此，这种方式使用的是HTTP短链接，每次请求都需要建立一个新的连接，请求结束后断开连接。
因此这种方式可以确保一应一答，且响应数据一定是返回给请求者。

痛点：但是如果是直接使用TCP长连接的服务器呢？
为了确保收到响应的一定是对应的请求者（可能是一个线程），需要唯一标识每个请求数据属于哪个请求者。
服务端还是按照之前的逻辑，只是在响应中原样抄上请求者ID。
然后再客户端收到服务端响应后，根据响应数据中的请求者ID，就可以路由到对应的请求者。

因此，这种方式要求：
1. 客户端需要能区分出各请求者
2. 客户端请求需要指定请求者ID
3. 客户端需要能路由到对应的请求者
*/

// go run demos/task/server/main.go demos/task/server/caculate.go
func main() {
	Log.SetLevel(LevelDebug)
	s := server.NewServer()
	factory := CalculateJobFactory{}
	s.Route(jobs.AddJobTag, factory.CreateCalculator(jobs.AddJobTag)).
		Route(jobs.SubJobTag, factory.CreateCalculator(jobs.SubJobTag)).
		Route(jobs.MulJobTag, factory.CreateCalculator(jobs.MulJobTag)).
		Route(jobs.DivJobTag, factory.CreateCalculator(jobs.DivJobTag))
	s.ListenAndServe()
	Log.Info("Server exit")
}
