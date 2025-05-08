package job

import iface "my-zinx/interface"

// 由于是channel，因此内部有锁，是并发安全的
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
