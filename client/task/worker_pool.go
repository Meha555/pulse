package task

import (
	"github.com/Meha555/pulse/core"
	"github.com/Meha555/pulse/utils"
)

type WorkerPool struct {
	core.WorkerPool[func()]
}

func NewWorkerPool(workers int, mq utils.IQueue[func()]) *WorkerPool {
	return &WorkerPool{
		WorkerPool: *core.NewWorkerPool(workers, mq, &TaskProcesser{}),
	}
}

type TaskProcesser struct{}

func (p *TaskProcesser) Process(fn func()) error {
	fn()
	return nil
}
