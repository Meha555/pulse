package utils

import (
	"sync"
	"sync/atomic"
)

type Dict[K comparable, V any] struct {
	// dict map[string]interface{}
	dict sync.Map
	size atomic.Int32 // TODO 近似值，因为没有把dict.Store和size.Store这种操作组合成原子操作。如果要做的话应该手动用无锁方式实现Dict类而不是靠sync.Map
	// mtx  sync.RWMutex
}

// func NewDict() *Dict {
// 	return &Dict{
// 		dict: make(map[string]interface{}),
// 	}
// }

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
	value, exists = val.(V)
	if !exists {
		panic(exists)
	}
	// value, exists = d.dict[key]
	return
}

func (d *Dict[K, V]) Store(key K, value V) {
	// d.mtx.Lock()
	// defer d.mtx.Unlock()
	// if d.dict == nil {
	// 	d.dict = make(map[string]interface{})
	// }
	d.dict.Store(key, value)
	d.size.Add(1)
	// d.dict[key] = value
}

func (d *Dict[K, V]) Delete(key K) {
	// d.mtx.Lock()
	// defer d.mtx.Unlock()
	// if d.dict == nil {
	// 	d.dict = make(map[string]interface{})
	// }
	if (d.size.Load() == 0) {
		return
	}
	d.dict.Delete(key)
	d.size.Add(-1)
	// delete(d.dict, key)
}

func (d *Dict[K, V]) Size() int32 {
	return d.size.Load()
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
