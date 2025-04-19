package job

import (
	iface "my-zinx/zinx/interface"
)

// 实现 job 时, 先嵌入这个基类, 然后根据需要对这个基类的方法进行重写
type BaseJob struct{}

/*
此处 BaseJob 的方法都为空的原因是有一些 Api 不需要 PreHandle 和 PostHandle
*/
func (b *BaseJob) PreHandle(req iface.IRequest) error  { return nil }
func (b *BaseJob) Handle(req iface.IRequest) error     { return nil }
func (b *BaseJob) PostHandle(req iface.IRequest) error { return nil }

type HeartBeatJob struct {
	BaseJob
}

func (h *HeartBeatJob) Handle(req iface.IRequest) error {
	// 心跳包（客户端主动发送，只有包头没有包体，服务端不回复）
	req.Conn().UpdateHeartBeat()
	return nil
}
