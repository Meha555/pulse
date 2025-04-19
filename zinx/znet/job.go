package znet

import (
	"my-zinx/zinx/ziface"
)

// 实现 job 时, 先嵌入这个基类, 然后根据需要对这个基类的方法进行重写
type BaseJob struct{}

/*
此处 BaseJob 的方法都为空的原因是有一些 Api 不需要 PreHandle 和 PostHandle
*/
func (b *BaseJob) PreHandle(req ziface.IRequest) error  { return nil }
func (b *BaseJob) Handle(req ziface.IRequest) error     { return nil }
func (b *BaseJob) PostHandle(req ziface.IRequest) error { return nil }

// TODO 结合context实现超时
// type HeartBeatApi struct {
// 	BaseJob

// 	seq uint32
// }

// func NewHeartBeatApi() *HeartBeatApi {
// 	return &HeartBeatApi{
// 		seq: rand.Uint32(),
// 	}
// }

// func (h *HeartBeatApi) Handle(req ziface.IRequest) error {
// 	conn := req.Conn().(*Connection)
// 	log.Println(req.Msg().(*SeqedTLVMsg).Body())

// 	return err
// }
