package server

type IServer interface {
	// 启动服务器，监听端口
	Listen()
	// 执行具体的服务器业务
	Serve()
	// 停止服务器
	Shutdown()
}
