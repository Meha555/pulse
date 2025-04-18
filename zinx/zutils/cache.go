package zutils

type Cache struct {
	dict map[string]string
}

func (e *Cache) Get(key string) (value string, exists bool) {
	if e.dict == nil {
		e.dict = make(map[string]string)
	}
	value, exists = e.dict[key]
	return
}

func (e *Cache) Set(key, value string) {
	if e.dict == nil {
		e.dict = make(map[string]string)
	}
	e.dict[key] = value
}
