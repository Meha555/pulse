package znet

import "my-zinx/zinx/ziface"

// 实现 router 时, 先嵌入这个基类, 然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

/*
此处 BaseRouter 的方法都为空的原因是有一些 Router 不需要 PreHandle 和 PostHandle
*/
func (br *BaseRouter) PreHandle(req ziface.IRequest)  {}
func (br *BaseRouter) Handle(req ziface.IRequest)     {}
func (br *BaseRouter) PostHandle(req ziface.IRequest) {}
