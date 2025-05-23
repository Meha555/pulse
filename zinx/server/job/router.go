package job

import (
	"fmt"

	"my-zinx/server/common"
	"my-zinx/utils"
)

// Api Tags
const (
	// 0-99是给用户预留的自定义tag

	HeartBeatTag = iota + 100
)

// IJobRouter
// Tag与路由的映射管理，根据Request中msg的Tag来确认用的是哪个路由，从而调用对应的3个回调
type IJobRouter interface {
	// 从管理器中获取指定的路由
	GetJob(tag uint16) IJob
	// 往管理器中添加一个路由
	AddJob(tag uint16, job IJob) IJobRouter
	// 执行指定的路由回调
	ExecJob(tag uint16, request common.IRequest) error
}

type JobRouter struct {
	// <tag, job>映射表
	apis *utils.Dict[uint16, IJob]
}

func NewJobRouter() *JobRouter {
	return &JobRouter{
		apis: utils.NewDict[uint16, IJob](),
	}
}

func (r *JobRouter) GetJob(tag uint16) IJob {
	job, ok := r.apis.Load(tag)
	if !ok {
		logger.Errorf("get job failed")
		return nil
	}
	return job
}

func (r *JobRouter) AddJob(tag uint16, job IJob) IJobRouter {
	r.apis.Store(tag, job)
	return r
}

func (r *JobRouter) ExecJob(tag uint16, req common.IRequest) error {
	if job, ok := r.apis.Load(tag); ok {
		if err := job.PreHandle(req); err != nil {
			return fmt.Errorf("call PreHandle error: %v", err)
		}
		if err := job.Handle(req); err != nil {
			return fmt.Errorf("call Handle error: %v", err)
		}
		if err := job.PostHandle(req); err != nil {
			return fmt.Errorf("call PostHandle error: %v", err)
		}
		return nil
	}
	return fmt.Errorf("no job for tag[%d]", tag)
}
