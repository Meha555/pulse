package task

import (
	"my-zinx/utils"
	"time"

	"github.com/google/uuid"
)

const (
	kDefaultCleanInterval = time.Second * 60 // 清理间隔s
	kDefaultThreshold     = 100
)

// taskTable 用于暂时存储任务信息
// 我们有时需要在任务结束后仍然暂存一下任务的信息，比如demo中的tasks就需要打印参数
// 但是又不能一直存着浪费内存，所以这里搞了一个类似于缓存的组件
type taskTable struct {
	utils.Dict[uuid.UUID, ITask]
	cleanInterval time.Duration

	cleanCn chan struct{}
}

type taskTableOption func(*taskTable)

func newTaskTable(opts ...taskTableOption) *taskTable {
	tbl := &taskTable{}

	for _, opt := range opts {
		opt(tbl)
	}

	if tbl.cleanInterval == 0 {
		tbl.cleanInterval = kDefaultCleanInterval
	}
	if tbl.Capacity() == 0 {
		tbl.SetCapacity(kDefaultThreshold)
	}

	// cleaner
	go func() {
		cleanFunc := func() {
			for val := range tbl.Iter() {
				id, task := val.Key, val.Value
				if task.CreateTime().Before(time.Now()) {
					if task.Status() != TaskStatusRunning && task.Status() != TaskStatusCreated {
						tbl.Delete(id)
					}
				} else {
					logger.Warnf("There is a task overstay for more than %v", tbl.cleanInterval)
				}
			}
		}

		ticker := time.NewTicker(tbl.cleanInterval)
		for {
			select {
			case <-tbl.cleanCn:
				cleanFunc()
				ticker.Reset(tbl.cleanInterval)
			case <-ticker.C:
				cleanFunc()
			}
		}
	}()

	return tbl
}

func WithCleanInterval(interval time.Duration) taskTableOption {
	return func(tbl *taskTable) {
		tbl.cleanInterval = interval
	}
}

func WithCapacity(cap int) taskTableOption {
	return func(tbl *taskTable) {
		tbl.SetCapacity(int32(cap))
	}
}

func (tbl *taskTable) Add(id uuid.UUID, task *Task) (err error) {
	if err = tbl.Store(id, task); err != nil {
		tbl.cleanCn <- struct{}{}
	}
	return
}

var TaskTbl *taskTable

func init() {
	TaskTbl = newTaskTable()
}
