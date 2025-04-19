package main

import (
	"fmt"
	"log"
	"my-zinx/zinx/ziface"
	"my-zinx/zinx/znet"
)

type EchoApi struct {
	znet.BaseJob
}

func (p *EchoApi) Handle(request ziface.IRequest) error {
	fmt.Println("Call Api Handle")
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
	s.ApiMapper.
		AddJob(0, &EchoApi{}).
		AddJob(1, &EchoApi{})
	s.ListenAndServe()
}
