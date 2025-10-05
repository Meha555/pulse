package server

// IBanner 展示启动欢迎信息
type IBanner interface {
	Show()
}

type IServer interface {
	// 启动服务器，监听端口
	Listen()
	// 执行具体的服务器业务
	Serve()
	// 停止服务器
	Shutdown()

	// 设置启动欢迎信息（因为允许不设置Banner，所以单独搞一个方法来注入，而不是在构造函数中）
	SetBanner(banner IBanner)
}
