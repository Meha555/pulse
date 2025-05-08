package iface

type IQueue[T any] interface {
	Push(request T)
	Pop() T
	Len() int
	Cap() int
	Close()
}