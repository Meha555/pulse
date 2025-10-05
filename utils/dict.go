package utils

import (
	"errors"
	"sync"
	"sync/atomic"
)

type IMap[K comparable, V any] interface {
	Store(key K, value any) error
	Load(key K) (value V, exists bool)
}

const (
	kDefaultCapacity = 256
)

// Dict 并发安全的字典类型
type Dict[K comparable, V any] struct {
	// dict map[string]interface{}
	dict sync.Map
	cap  atomic.Int32
	size atomic.Int32 // TODO 近似值，因为没有把dict.Store和size.Store这种操作组合成原子操作。如果要做的话应该手动用无锁方式实现Dict类而不是靠sync.Map
	// mtx  sync.RWMutex
}

var ErrDictIsFull = errors.New("dict is full")

type dictOption[K comparable, V any] func(*Dict[K, V])

func WithCapacity[K comparable, V any](cap int) dictOption[K, V] {
	return func(d *Dict[K, V]) {
		d.SetCapacity(int32(cap))
	}
}

func NewDict[K comparable, V any](opts ...dictOption[K, V]) *Dict[K, V] {
	d := &Dict[K, V]{
		cap: atomic.Int32{},
	}

	for _, opt := range opts {
		opt(d)
	}

	d.cap.CompareAndSwap(0, kDefaultCapacity)

	return d
}

func (d *Dict[K, V]) Load(key K) (value V, exists bool) {
	// d.mtx.RLock()
	// defer d.mtx.RUnlock()
	// if d.dict == nil {
	// 	d.dict = make(map[string]interface{})
	// }
	val, exists := d.dict.Load(key)
	if !exists {
		return
	}
	value = val.(V)
	// value, exists = d.dict[key]
	return
}

func (d *Dict[K, V]) Store(key K, value V) error {
	// 如果 CompareAndSwap 失败，说明 size 被其他 goroutine 修改，继续循环重试
	for {
		// d.mtx.Lock()
		// defer d.mtx.Unlock()
		// if d.dict == nil {
		// 	d.dict = make(map[string]interface{})
		// }
		oldSize := d.size.Load()
		if oldSize >= d.cap.Load() {
			return ErrDictIsFull
		}
		if d.size.CompareAndSwap(oldSize, oldSize+1) {
			d.dict.Store(key, value)
			return nil
		}
		// d.dict[key] = value
	}
}

func (d *Dict[K, V]) Delete(key K) {
	// d.mtx.Lock()
	// defer d.mtx.Unlock()
	// if d.dict == nil {
	// 	d.dict = make(map[string]interface{})
	// }
	if d.size.CompareAndSwap(0, 0) {
		return
	} else {
		d.dict.Delete(key)
		d.size.Add(-1)
	}
	// delete(d.dict, key)
}

func (d *Dict[K, V]) Size() int32 {
	return d.size.Load()
}

func (d *Dict[K, V]) Capacity() int32 {
	return d.cap.Load()
}

func (d *Dict[K, V]) SetCapacity(cap int32) {
	d.cap.Store(cap)
}

func (d *Dict[K, V]) Range(f func(key K, value V) bool) {
	d.dict.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}

// Iter 返回一个只读通道，用于迭代所有键值对
func (d *Dict[K, V]) Iter() <-chan struct {
	Key   K
	Value V
} {
	ch := make(chan struct {
		Key   K
		Value V
	})

	go func() {
		d.dict.Range(func(key, value interface{}) bool {
			ch <- struct {
				Key   K
				Value V
			}{
				Key:   key.(K),
				Value: value.(V),
			}
			return true
		})
		close(ch)
	}()

	return ch
}

// 由于 comparable 不能在类型约束外使用，这里需要使用具体的类型参数。
// 假设 IMap 接口期望的是具体的可比较类型，这里使用 string 作为示例。
var _ IMap[string, any] = (*Dict[string, any])(nil)
