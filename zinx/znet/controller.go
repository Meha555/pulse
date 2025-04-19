package znet

import (
	"my-zinx/zinx/ziface"
)

// 实现 controller 时, 先嵌入这个基类, 然后根据需要对这个基类的方法进行重写
type BaseController struct{}

/*
此处 BaseController 的方法都为空的原因是有一些 Controller 不需要 PreHandle 和 PostHandle
*/
func (b *BaseController) PreHandle(req ziface.IRequest) error  { return nil }
func (b *BaseController) Handle(req ziface.IRequest) error     { return nil }
func (b *BaseController) PostHandle(req ziface.IRequest) error { return nil }

// TODO 结合context实现超时
// type HeartBeatController struct {
// 	BaseController

// 	seq uint32
// }

// func NewHeartBeatController() *HeartBeatController {
// 	return &HeartBeatController{
// 		seq: rand.Uint32(),
// 	}
// }

// func (h *HeartBeatController) Handle(req ziface.IRequest) error {
// 	conn := req.Conn().(*Connection)
// 	log.Println(req.Msg().(*SeqedTLVMsg).Body())

// 	return err
// }
