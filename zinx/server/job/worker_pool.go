package job

import (
	"my-zinx/core"
	"my-zinx/server/common"
	"my-zinx/utils"
)

type MsgQueue = utils.BlockingQueue[common.IRequest]

// WorkerPool 业务协程池
// - 支持启动多个业务协程，每个业务协程执行一个业务操作。
// - 支持等待业务完成。
// - 支持提交消息到消息队列，并由空闲协程执行
type WorkerPool struct {
	core.WorkerPool[common.IRequest]
	router IJobRouter // API映射（引用）
}

func NewWorkerPool(workers int, mq utils.IQueue[common.IRequest], router IJobRouter) *WorkerPool {
	return &WorkerPool{
		WorkerPool: *core.NewWorkerPool(workers, mq, &JobProcesser{router: router}),
		router:     router,
	}
}

type JobProcesser struct {
	router IJobRouter // API映射（引用）
}

func (p *JobProcesser) Process(request common.IRequest) error {
	return p.router.ExecJob(request.Msg().Tag(), request)
}
