package job

import (
	iface "my-zinx/interface"
	"my-zinx/utils"
	"sync"
)

type MsgQueue = utils.BlockingQueue[iface.IRequest]

// WorkerPool 业务协程池
// - 支持启动多个业务协程，每个业务协程执行一个业务操作。
// - 支持等待业务完成。
// - 支持提交消息到消息队列，并由空闲协程执行
type WorkerPool struct {
	workers int                          // 协程池大小
	mq      iface.IQueue[iface.IRequest] // 消息队列（引用）
	router  iface.IJobRouter             // API映射（引用）
	stopCh  chan struct{}                // 停止信号通道
	wg      sync.WaitGroup
}

func NewWorkerPool(workers int, mq iface.IQueue[iface.IRequest], router iface.IJobRouter) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		mq:      mq,
		router:  router,
		stopCh:  make(chan struct{}),
	}
}

// Start 初始化协程池，启动指定数量的协程并阻塞等待mq中的消息。
func (w *WorkerPool) Start() {
	for i := range w.workers {
		w.wg.Add(1)
		go func(workerID int) {
			defer w.wg.Done()
			logger.Debugf("Worker[%d] started", workerID)
			for {
				select {
				case <-w.stopCh:
					logger.Debugf("Worker[%d] stopping", workerID)
					return
				default:
					if request := w.mq.Pop(); request != nil {
						logger.Debugf("Worker[%d] processing request", workerID)
						if err := w.processRequest(request); err != nil {
							logger.Errorf("Worker[%d] process request failed: %v", workerID, err)
						}
					}
				}
			}
		}(i)
	}
}

// Stop 停止协程池，等待所有协程退出
func (w *WorkerPool) Stop() {
	w.mq.Close()
	close(w.stopCh) // 发送停止信号
	w.wg.Wait()     // 等待所有协程退出
	logger.Debug("All workers stopped")
}

func (w *WorkerPool) Post(request iface.IRequest) {
	if request == nil {
		return
	}
	w.mq.Push(request)
}

// processRequest 执行具体的业务逻辑
func (w *WorkerPool) processRequest(request iface.IRequest) error {
	return w.router.ExecJob(request.Msg().Tag(), request)
}
