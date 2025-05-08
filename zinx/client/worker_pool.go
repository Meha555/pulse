package client

import (
	iface "my-zinx/interface"
	"my-zinx/utils"
	"sync"
)

type TaskQueue = utils.BlockingQueue[TaskHandler]

type WorkerPool struct {
	workers int
	mq      iface.IQueue[TaskHandler]
	stopCh  chan struct{}
	wg      sync.WaitGroup
}

func NewWorkerPool(workers int, mq iface.IQueue[TaskHandler]) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		mq:      mq,
		stopCh:  make(chan struct{}),
	}
}

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
					if handler := w.mq.Pop(); handler != nil {
						logger.Debugf("Worker[%d] processing request", workerID)
						if err := handler.Do(); err != nil {
							logger.Errorf("Worker[%d] process request failed: %v", workerID, err)
						}
					}
				}
			}
		}(i)
	}
}

func (w *WorkerPool) Stop() {
	w.mq.Close()
	close(w.stopCh) // 发送停止信号
	w.wg.Wait()     // 等待所有协程退出
	logger.Debug("All workers stopped")
}

func (w *WorkerPool) Post(handler TaskHandler) {
	if handler == nil {
		return
	}
	w.mq.Push(handler)
}
