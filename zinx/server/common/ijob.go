package common

// IJob 具体的业务逻辑
type IJob interface {
	// 处理业务之前的钩子方法
	PreHandle(req IRequest) error
	// 处理业务的主方法
	Handle(req IRequest) error
	// 处理业务之后的钩子方法
	PostHandle(req IRequest) error
}
