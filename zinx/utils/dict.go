package utils

import "sync"

type Dict struct {
	dict map[string]interface{}
	mtx  sync.RWMutex
}

func NewDict() *Dict {
	return &Dict{
		dict: make(map[string]interface{}),
	}
}

func (d *Dict) Get(key string) (value interface{}, exists bool) {
	d.mtx.RLock()
	defer d.mtx.RUnlock()
	if d.dict == nil {
		d.dict = make(map[string]interface{})
	}
	value, exists = d.dict[key]
	return
}

func (d *Dict) Set(key string, value interface{}) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	if d.dict == nil {
		d.dict = make(map[string]interface{})
	}
	d.dict[key] = value
}

func (d *Dict) Del(key string) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	if d.dict == nil {
		d.dict = make(map[string]interface{})
	}
	delete(d.dict, key)
}
