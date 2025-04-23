package main

import (
	"context"
	"fmt"
	"my-zinx/zinx/core/job"
	"my-zinx/zinx/core/message"
	"my-zinx/zinx/core/session"
	"net"
	"time"
)

func main() {
	fmt.Println("Client Test ... start")
	ep, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:3333")
	peer, err := net.DialTCP("tcp4", nil, ep)
	if err != nil {
		fmt.Println("client start err: ", err)
		return
	}
	conn := session.NewConnection(peer, context.Background(), nil)

	go doHeartBeat(conn)
	go doEcho(conn, 0)
	go doEcho(conn, 1)

	select {}
}

func doHeartBeat(conn *session.Session) {
	ticker := time.NewTicker(time.Second)
	for curTime := range ticker.C {
		msgSent := message.NewSeqedTLVMsg(0, job.HeartBeatTag, nil)
		data, err := message.Marshal(msgSent)
		if err != nil {
			fmt.Println("Marshal error:", err)
			continue
		}
		_, err = conn.Send(data)
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}
		fmt.Println("heartbeat success: ", curTime)
	}
}

func doEcho(conn *session.Session, id uint16) {
	var serial uint32 = 0
	for {
		msgSent := message.NewSeqedTLVMsg(serial, id, fmt.Appendf(nil, "hello ZINX %d", id))
		data, err := message.Marshal(msgSent)
		if err != nil {
			fmt.Println("Marshal error:", err)
			continue
		}
		_, err = conn.Send(data)
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}
		serial++

		msg := &message.SeqedTLVMsg{}
		err = conn.RecvMsg(msg)
		if err != nil {
			fmt.Println("read buf error:", err)
			return
		}

		fmt.Printf("read: %d %s\n", msg.Serial(), string(msg.Body()))

		time.Sleep(time.Duration(1<<id) * time.Second)
	}
}
