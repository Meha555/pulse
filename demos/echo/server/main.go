package main

import (
	"fmt"
	"log"
	"my-zinx/zinx/core/connection"
	"my-zinx/zinx/core/job"
	iface "my-zinx/zinx/interface"
	"my-zinx/zinx/server"
)

type EchoJob struct {
	job.BaseJob
}

func (p *EchoJob) Handle(request iface.IRequest) error {
	fmt.Println("Call Api Handle")
	msg := request.Msg()
	log.Printf("ReadMsg: %d %s\n", msg.Serial(), string(msg.Body()))
	if err := request.Conn().(*connection.Connection).SendMsg(msg); err != nil {
		log.Println("Write error:", err)
		return err
	} else {
		log.Println("Write success")
		return nil
	}
}

func main() {
	s := server.NewServer()
	s.JobRouter.
		AddJob(0, &EchoJob{}).
		AddJob(1, &EchoJob{})
	s.ListenAndServe()
	fmt.Println("Server exit")
}
