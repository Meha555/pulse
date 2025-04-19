package utils

import (
	"context"
)

type Context struct {
	context.Context // TODO
	dict            *Dict
}

func NewContext(parent context.Context) Context {
	return Context{
		parent, NewDict(),
	}
}

func (c *Context) Set(key string, value interface{}) {
	c.dict.Set(key, value)
}

func (c *Context) Get(key string) (interface{}, bool) {
	return c.dict.Get(key)
}

func (c *Context) Del(key string) {
	c.dict.Del(key)
}
