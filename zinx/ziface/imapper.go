package ziface

// IControllerMapper
// Tag与路由的映射管理
type IControllerMapper interface {
	// 从管理器中获取指定的路由
	GetController(tag uint16) IController
	// 往管理器中添加一个路由
	AddController(tag uint16, controller IController) IControllerMapper
	// 执行指定的路由回调
	ExecController(tag uint16, request IRequest) error
}
