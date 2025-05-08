package main

import (
	"fmt"
	"my-zinx/zinx/core/job"
	"my-zinx/zinx/core/session"
	iface "my-zinx/zinx/interface"
	. "my-zinx/zinx/log"
	"my-zinx/zinx/server"
)

type EchoJob struct {
	job.BaseJob
}

func (p *EchoJob) Handle(req iface.IRequest) error {
	fmt.Println("Call Api Handle")
	msg := req.Msg()
	Log.Infof("ReadMsg: %d %s\n", msg.Serial(), string(msg.Body()))
	if err := req.Session().(*session.Session).SendMsg(msg); err != nil {
		Log.Errorf("Write error: %v", err)
		return err
	} else {
		Log.Info("Write success")
		return nil
	}
}

func main() {
	s := server.NewServer()
	s.Route(0, &EchoJob{}).Route(1, &EchoJob{})
	s.ListenAndServe()
	Log.Info("Server exit")
}
