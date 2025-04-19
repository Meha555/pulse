package znet

import "my-zinx/zinx/ziface"

type MsgQueue chan ziface.IRequest

func NewMsgQueue(capacity int) *MsgQueue {
	queue := make(MsgQueue, capacity)
	return &queue
}

func (t MsgQueue) Push(request ziface.IRequest) {
	t <- request
}

func (t MsgQueue) Pop() ziface.IRequest {
	return <-t
}

func (t MsgQueue) Len() int {
	return len(t)
}

func (t MsgQueue) Cap() int {
	return cap(t)
}

var _ ziface.IQueue = (*MsgQueue)(nil)
