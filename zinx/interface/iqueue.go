package iface

type IQueue interface {
	Push(request IRequest)
	Pop() IRequest
	Len() int
	Cap() int
	Close()
}
