package core

import (
	"pulse/logging"
	"pulse/utils"
	"sync"
)

var logger = logging.NewStdLogger(logging.LevelInfo, "core", "[%t] [%c %l] [%f:%C:%L:%g] %m", false)

type Processer[Handler any] interface {
	// Process 执行处理逻辑。需要实现者在其中处理panic，否则协程会退出
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
	wp := &WorkerPool[Handler]{
		workers:   workers,
		mq:        mq,
		stopCh:    make(chan struct{}),
		processer: processer,
	}

	return wp
}

func (w *WorkerPool[Hanlder]) Start() {
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
					logger.Debugf("Worker[%d] processing request", workerID)
					handler := w.mq.Pop()
					if err := w.processer.Process(handler); err != nil {
						logger.Errorf("Worker[%d] process request failed: %v", workerID, err)
					}
				}
			}
		}(i)
	}
}

func (w *WorkerPool[Handler]) Stop() {
	w.mq.Close()
	close(w.stopCh)
	w.wg.Wait()
	logger.Debug("All workers stopped")
}

func (w *WorkerPool[Handler]) Post(handler Handler) {
	w.mq.Push(handler)
}
