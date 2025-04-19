package ziface

// IController
// 路由抽象接口，路由里的数据都是Request。根据Request中msg的Tag来确认用的是哪个路由，从而调用对应的3个回调
// TODO 我想的是允许每步对请求进行修改，是否要将入参改为指针呢？但是这样的话实现接口就不对了
type IController interface {
	// 处理业务之前的钩子方法
	PreHandle(request IRequest) error
	// 处理业务的主方法
	Handle(request IRequest) error
	// 处理业务之后的钩子方法
	PostHandle(request IRequest) error
}
