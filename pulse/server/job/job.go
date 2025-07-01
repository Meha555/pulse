package job

import (
	"pulse/logging"
	"pulse/server/common"
)

var logger = logging.NewStdLogger(logging.LevelInfo, "job", "[%t] [%c %l] [%f:%C:%L:%g] %m", false)

// IJob 具体的业务逻辑
type IJob interface {
	// 处理业务之前的钩子方法
	PreHandle(req common.IRequest) error
	// 处理业务的主方法
	Handle(req common.IRequest) error
	// 处理业务之后的钩子方法
	PostHandle(req common.IRequest) error
}

// 实现 job 时, 先嵌入这个基类, 然后根据需要对这个基类的方法进行重写
type BaseJob struct{}

/*
此处 BaseJob 的方法都为空的原因是有一些 Api 不需要 PreHandle 和 PostHandle
*/
func (b *BaseJob) PreHandle(req common.IRequest) error  { return nil }
func (b *BaseJob) Handle(req common.IRequest) error     { return nil }
func (b *BaseJob) PostHandle(req common.IRequest) error { return nil }

type HeartBeatJob struct {
	BaseJob
}

func (h *HeartBeatJob) Handle(req common.IRequest) error {
	// 心跳包（客户端主动发送，只有包头没有包体）
	req.Session().UpdateHeartBeat()
	return nil
}
