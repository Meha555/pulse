package ziface

// IApiMapper
// Tag与路由的映射管理，根据Request中msg的Tag来确认用的是哪个路由，从而调用对应的3个回调
type IApiMapper interface {
	// 从管理器中获取指定的路由
	GetJob(tag uint16) IJob
	// 往管理器中添加一个路由
	AddJob(tag uint16, job IJob) IApiMapper
	// 执行指定的路由回调
	ExecJob(tag uint16, request IRequest) error
}
