package main

import (
	"context"
	"fmt"
	"my-zinx/zinx/core/job"
	"my-zinx/zinx/core/message"
	"my-zinx/zinx/core/session"
	. "my-zinx/zinx/log"
	"net"
	"time"
)

func main() {
	Log.Info("Client Test ... start")
	ep, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:3333")
	peer, err := net.DialTCP("tcp4", nil, ep)
	if err != nil {
		Log.Errorf("client start err: %v", err)
		return
	}
	conn := session.NewSession(peer, context.Background(), nil)

	go doHeartBeat(conn)
	go doEcho(conn, 1)
	go doEcho(conn, 2)

	select {}
}

func doHeartBeat(conn *session.Session) {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		msgSent := message.NewSeqedTLVMsg(0, job.HeartBeatTag, nil)
		data, err := message.Marshal(msgSent)
		if err != nil {
			Log.Errorf("Marshal error: %v", err)
			continue
		}
		_, err = conn.Send(data)
		if err != nil {
			Log.Errorf("Write error: %v", err)
			return
		}
	}
}

func doEcho(conn *session.Session, id uint16) {
	var serial uint32 = 0
	for {
		msgSent := message.NewSeqedTLVMsg(serial, id, fmt.Appendf(nil, "hello ZINX %d", id))
		data, err := message.Marshal(msgSent)
		if err != nil {
			Log.Errorf("Marshal error: %v", err)
			continue
		}
		_, err = conn.Send(data)
		if err != nil {
			Log.Errorf("Write error: %v", err)
			return
		}
		serial++

		msg := &message.SeqedTLVMsg{}
		err = conn.RecvMsg(msg)
		if err != nil {
			Log.Errorf("read buf error: %v", err)
			return
		}

		Log.Infof("read: %d %d %s\n", msg.Serial(), msg.Tag(), string(msg.Body()))

		time.Sleep(time.Duration(1<<id) * time.Second)
	}
}
