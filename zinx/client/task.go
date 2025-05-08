package client

import (
	"sync"

	"github.com/google/uuid"
)

type ITask interface {
	ID() uuid.UUID
	Exec() error
	Data() []interface{}
	AppendData(userData ...interface{})
}

type TaskHandler interface {
	Do() error
}

type TaskHandlerFunc func() error

func (t TaskHandlerFunc) Do() error {
	return t()
}

type Task struct {
	id     uuid.UUID
	fn     TaskHandlerFunc
	data   []interface{}
	repeat int

	pool *WorkerPool
}

type option func(*Task)

func WithWorkerPool(pool *WorkerPool) option {
	return func(t *Task) {
		t.pool = pool
	}
}

func WithRepeat(repeat int) option {
	return func(t *Task) {
		t.repeat = repeat
	}
}

func WithData(userData ...interface{}) option {
	return func(t *Task) {
		t.data = userData
	}
}

// NewTask 创建一个异步任务
// 在fn内部可以直接通过闭包的方式访问到Task实例。设计思路参考testing.T
func NewTask(id uuid.UUID, fn func(*Task) error, opts ...option) *Task {
	t := &Task{
		id: id,
	}
	t.fn = func() error {
		return fn(t)
	}

	for _, opt := range opts {
		opt(t)
	}

	mu.Lock()
	taskMap[t.id] = t
	mu.Unlock()
	return t
}

func (t *Task) ID() uuid.UUID {
	return t.id
}

func (t *Task) Exec() error {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Task %s panic: %v", t.id, err)
		}
		if t.repeat > 0 {
			t.repeat--
			if t.pool != nil {
				t.pool.Post(t.fn)
			} else {
				t.Exec()
			}
			return
		}
		mu.Lock()
		delete(taskMap, t.id)
		mu.Unlock()
	}()
	if t.pool != nil {
		t.pool.Post(t.fn)
		return nil
	} else {
		return t.fn()
	}
}

func (t *Task) Data() []interface{} {
	return t.data
}

func (t *Task) AppendData(userData ...interface{}) {
	t.data = append(t.data, userData...)
}

var (
	taskMap = make(map[uuid.UUID]ITask)
	mu      sync.Mutex
)

func GetTask(id uuid.UUID) ITask {
	mu.Lock()
	defer mu.Unlock()
	return taskMap[id]
}
