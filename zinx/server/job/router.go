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

type JobRouter struct {
	// <tag, job>映射表
	apis utils.Dict[uint16, common.IJob]
}

func NewJobRouter() *JobRouter {
	return &JobRouter{}
}

func (r *JobRouter) GetJob(tag uint16) common.IJob {
	job, ok := r.apis.Load(tag)
	if !ok {
		logger.Errorf("get job failed")
		return nil
	}
	return job
}

func (r *JobRouter) AddJob(tag uint16, job common.IJob) common.IJobRouter {
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
