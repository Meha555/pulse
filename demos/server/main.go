package main

import (
	"fmt"
	"log"
	"my-zinx/zinx/ziface"
	"my-zinx/zinx/znet"
)

type EchoRouter struct {
	znet.BaseRouter
}

func (p *EchoRouter) PreHandle(request ziface.IRequest) error {
	fmt.Println("Call Router PreHandle")
	return nil
}

func (p *EchoRouter) Handle(request ziface.IRequest) error {
	fmt.Println("Call Router Handle")
	if nbytes, err := request.Conn().(*znet.Connection).Send(request.Data()); err != nil {
		log.Println("Write error:", err)
		return err
	} else {
		log.Println("Write success, nbytes:", nbytes)
		return nil
	}
}

func (p *EchoRouter) PostHandle(request ziface.IRequest) error {
	fmt.Println("Call Router PostHandle")
	return nil
}

func main() {
	s := znet.NewServer("[ZINX v0.3]", 3333)
	s.AddRouter(&EchoRouter{})
	s.Start()
	s.Serve()
}
