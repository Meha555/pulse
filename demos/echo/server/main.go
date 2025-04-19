package main

import (
	"fmt"
	"log"
	"my-zinx/zinx/ziface"
	"my-zinx/zinx/znet"
)

type EchoController struct {
	znet.BaseController
}

func (p *EchoController) Handle(request ziface.IRequest) error {
	fmt.Println("Call Controller Handle")
	msg := request.Msg()
	log.Printf("ReadMsg: %d %s\n", msg.Serial(), string(msg.Body()))
	if err := request.Conn().(*znet.Connection).SendMsg(msg); err != nil {
		log.Println("Write error:", err)
		return err
	} else {
		log.Println("Write success")
		return nil
	}
}

func main() {
	s := znet.NewServer()
	s.ControllerMapper.
		AddController(0, &EchoController{}).
		AddController(1, &EchoController{})
	s.ListenAndServe()
}
