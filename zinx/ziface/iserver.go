package ziface

type IServer interface {
	// 启动服务器，监听端口
	Start()
	// 执行具体的服务器业务
	Serve()
	// 停止服务器
	Stop()
	// 添加一个路由对象来执行路由业务方法
	AddRouter(router IRouter)
}

