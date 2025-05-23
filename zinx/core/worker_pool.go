package core

import (
	"my-zinx/log"
	"my-zinx/utils"
	"sync"
)

var logger = log.NewStdLogger(log.LevelInfo, "core", "[%t] [%c %l] [%f:%C:%L:%g] %m", false)

type Processer[Handler any] interface {
	Process(Handler) error
}

type WorkerPool[Handler any] struct {
	workers int
	mq      utils.IQueue[Handler]
	stopCh  chan struct{}
	wg      sync.WaitGroup

	processer Processer[Handler]
}

func NewWorkerPool[Handler any](workers int, mq utils.IQueue[Handler], processer Processer[Handler]) *WorkerPool[Handler] {
	return &WorkerPool[Handler]{
		workers:   workers,
		mq:        mq,
		stopCh:    make(chan struct{}),
		processer: processer,
	}
}

func (w *WorkerPool[T]) Start() {
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
					var zero T
					if handler := w.mq.Pop(); any(handler) != any(zero) {
						logger.Debugf("Worker[%d] processing request", workerID)
						if err := w.processer.Process(handler); err != nil {
							logger.Errorf("Worker[%d] process request failed: %v", workerID, err)
						}
					}
				}
			}
		}(i)
	}
}

func (w *WorkerPool[T]) Stop() {
	w.mq.Close()
	close(w.stopCh)
	w.wg.Wait()
	logger.Debug("All workers stopped")
}

func (w *WorkerPool[T]) Post(handler T) {
	// 由于泛型类型 T 不能直接与 nil 比较，这里使用类型断言来处理
	var zero T
	if any(handler) == any(zero) {
		return
	}
	w.mq.Push(handler)
}
