package message

import iface "my-zinx/zinx/interface"

type MsgQueue chan iface.IRequest

func NewMsgQueue(capacity int) *MsgQueue {
	queue := make(MsgQueue, capacity)
	return &queue
}

func (t MsgQueue) Push(request iface.IRequest) {
	t <- request
}

func (t MsgQueue) Pop() iface.IRequest {
	return <-t
}

func (t MsgQueue) Len() int {
	return len(t)
}

func (t MsgQueue) Cap() int {
	return cap(t)
}

func (t MsgQueue) Close() {
	close(t)
}

var _ iface.IQueue = (*MsgQueue)(nil)
