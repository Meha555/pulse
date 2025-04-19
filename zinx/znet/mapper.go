package znet

import (
	"fmt"
	"my-zinx/zinx/ziface"
)

// Api Tags
const (
	// 0-99是给用户预留的自定义tag

	TAG_HEARTBEAT = iota + 100
)

type ApiMapper struct {
	// <tag, job>映射表
	apis map[uint16]ziface.IJob
}

func NewApiMapper() *ApiMapper {
	return &ApiMapper{
		apis: make(map[uint16]ziface.IJob),
	}
}

func (r *ApiMapper) GetJob(tag uint16) ziface.IJob {
	return r.apis[tag]
}

func (r *ApiMapper) AddJob(tag uint16, job ziface.IJob) ziface.IApiMapper {
	r.apis[tag] = job
	return r
}

func (r *ApiMapper) ExecJob(tag uint16, request ziface.IRequest) error {
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

var _ ziface.IApiMapper = (*ApiMapper)(nil)
