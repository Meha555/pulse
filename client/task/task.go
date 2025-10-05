package task

import (
	"context"
	"fmt"
	"time"

	"github.com/Meha555/go-tinylog"
	"github.com/google/uuid"
)

var logger *tinylog.Logger

func init() {
	var err error
	logger, err = tinylog.NewStdLogger(tinylog.LevelInfo, "task", "[%t] [%c %l] [%f:%C:%L:%g] %m", false, tinylog.Lcolored)
	if err != nil {
		panic(err)
	}
}

type Status int

func (s Status) String() string {
	switch s {
	case TaskStatusCreated:
		return "created"
	case TaskStatusPending:
		return "pending"
	case TaskStatusRunning:
		return "running"
	case TaskStatusFinished:
		return "finished"
	case TaskStatusFailed:
		return "failed"
	case TaskStatusCanceled:
		return "canceled"
	default:
		return "unknown"
	}
}

// 任务状态
const (
	TaskStatusCreated Status = iota
	TaskStatusPending
	TaskStatusRunning
	TaskStatusFinished
	TaskStatusFailed
	TaskStatusCanceled
)

// ITask 任务接口
// 参考了context，但接口语义不完全一致
type ITask interface {
	ID() uuid.UUID
	Exec()
	// 获取user data
	Data() []interface{}
	AppendData(userData ...interface{})
	Status() Status
	Cancel()
	// 创建时间
	CreateTime() time.Time
	// 返回超时时间，超时后会自动取消。没有超时则返回零值事件和false
	Deadline() (deadline time.Time, ok bool)
	// 检查任务是否完成/取消/失败
	Done() <-chan struct{}
	// 返回最后的错误
	Err() error
}

type TaskHandler interface {
	Do(*Task) error
}

type TaskHandlerFunc func(*Task) error

func (f TaskHandlerFunc) Do(t *Task) error {
	return f(t)
}

type Task struct {
	id uuid.UUID
	fn func()
	// data非并发安全，调用Exec后就不应该再修改
	data []interface{}
	// repeat执行的任务不之间可能并发，因为是上一个fn执行完后才会将下一个fn加入调度，这个特性很重要，避免了加锁
	repeat     int
	status     Status
	createTime time.Time
	ctx        context.Context
	cancel     context.CancelFunc
	doneCh     chan struct{}
	lastErr    error
	pool       *WorkerPool
}

type option func(*Task)

func WithWorkerPool(pool *WorkerPool) option {
	return func(t *Task) {
		t.pool = pool
	}
}

func WithRepeat(repeat int) option {
	return func(t *Task) {
		t.repeat = repeat - 1
	}
}

// 设置总共调用的超时（即使是repeat，最后一次必须在deadline之前）
// 因为这种既有超时也有重复的场景往往是：意图在有限时间内多次重试
func WithTimeout(timeout time.Duration) option {
	return func(t *Task) {
		t.ctx, t.cancel = context.WithTimeout(context.Background(), timeout)
	}
}

func WithData(userData ...interface{}) option {
	return func(t *Task) {
		t.data = userData
	}
}

// NewTask 创建一个异步任务（有协程池用协程池，没有的话就就地起一个协程）
// 在fn内部可以直接通过闭包的方式（以自身的指针作为参数）访问到Task实例。设计思路参考testing.T
func NewTask(id uuid.UUID, fn TaskHandlerFunc, opts ...option) *Task {
	t := &Task{
		id:         id,
		status:     TaskStatusCreated,
		doneCh:     make(chan struct{}, 1),
		createTime: time.Now(),
	}

	for _, opt := range opts {
		opt(t)
	}

	// 提供一个闭包函数
	t.fn = func() {
		defer func() {
			if t.cancel != nil {
				t.cancel()
			}
			if err := recover(); err != nil {
				logger.Errorf("Task %s panic: %v", t.id, err)
				t.status = TaskStatusFailed
				t.lastErr = fmt.Errorf("task %s panic: %v", t.id, err)
				close(t.doneCh)
			}
		}()
		if t.ctx != nil {
			if deadline, ok := t.ctx.Deadline(); ok && time.Now().After(deadline) {
				t.status = TaskStatusCanceled
				return
			}
		}
		if t.status == TaskStatusCanceled {
			t.doneCh <- struct{}{}
			return
		}
		t.lastErr = fn.Do(t)
		if t.repeat > 0 && t.status == TaskStatusRunning {
			t.repeat--
			t.Exec()
			<-t.doneCh
		}
		t.status = TaskStatusFinished
		close(t.doneCh)
	}

	TaskTbl.Store(t.id, t)
	return t
}

func (t *Task) ID() uuid.UUID {
	return t.id
}

func (t *Task) Status() Status {
	return t.status
}

func (t *Task) Err() error {
	return t.lastErr
}

func (t *Task) Exec() {
	if t.ctx != nil {
		if deadline, ok := t.ctx.Deadline(); ok && time.Now().After(deadline) {
			t.status = TaskStatusCanceled
		}
	}
	if t.status == TaskStatusCanceled || t.status == TaskStatusFinished || t.status == TaskStatusFailed {
		t.doneCh <- struct{}{}
		return
	}
	t.status = TaskStatusPending
	if t.pool != nil {
		t.pool.Post(t.fn)
	} else {
		t.status = TaskStatusRunning
		go t.fn()
	}
}

func (t *Task) Data() []interface{} {
	return t.data
}

func (t *Task) AppendData(userData ...interface{}) {
	t.data = append(t.data, userData...)
}

func (t *Task) Cancel() {
	t.status = TaskStatusCanceled
}

func (t *Task) CreateTime() time.Time {
	return t.createTime
}

func (t *Task) Deadline() (deadline time.Time, ok bool) {
	return t.ctx.Deadline()
}

func (t *Task) Done() <-chan struct{} {
	return t.doneCh
}
