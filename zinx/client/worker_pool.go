package client

import (
	"my-zinx/core"
	"my-zinx/utils"
)

type WorkerPool struct {
	core.WorkerPool[TaskHandler]
}

func NewWorkerPool(workers int, mq utils.IQueue[TaskHandler]) *WorkerPool {
	return &WorkerPool{
		WorkerPool: *core.NewWorkerPool(workers, mq, &TaskProcesser{}),
	}
}

type TaskProcesser struct{}

func (p *TaskProcesser) Process(handler TaskHandler) error {
	return handler.Do()
}
