package job

import (
	"fmt"
	iface "my-zinx/zinx/interface"
)

// Api Tags
const (
	// 0-99是给用户预留的自定义tag

	HeartBeatTag = iota + 100
)

type JobRouter struct {
	// <tag, job>映射表
	apis map[uint16]iface.IJob
}

func NewJobRouter() *JobRouter {
	return &JobRouter{
		apis: make(map[uint16]iface.IJob),
	}
}

func (r *JobRouter) GetJob(tag uint16) iface.IJob {
	return r.apis[tag]
}

func (r *JobRouter) AddJob(tag uint16, job iface.IJob) iface.IJobRouter {
	r.apis[tag] = job
	return r
}

func (r *JobRouter) ExecJob(tag uint16, request iface.IRequest) error {
	if job, ok := r.apis[tag]; ok {
		if err := job.PreHandle(request); err != nil {
			return fmt.Errorf("PreHandle error: %v", err)
		}
		if err := job.Handle(request); err != nil {
			return fmt.Errorf("Handle error: %v", err)
		}
		if err := job.PostHandle(request); err != nil {
			return fmt.Errorf("PostHandle error: %v", err)
		}
		return nil
	}
	return fmt.Errorf("no job for tag[%d]", tag)
}

var _ iface.IJobRouter = (*JobRouter)(nil)
