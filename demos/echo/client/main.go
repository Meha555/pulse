package main

import (
	"fmt"
	"my-zinx/zinx/znet"
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
	conn := znet.NewConnection(peer, nil)

	go doEcho(conn, 0)
	go doEcho(conn, 1)

	select {}
}

func doEcho(conn *znet.Connection, id int) {
	var serial uint32 = 0
	for {
		msgSent := znet.NewSeqedTLVMsg(serial, 0, fmt.Appendf(nil, "hello ZINX %d", id))
		data, err := znet.Marshal(msgSent)
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

		msg := &znet.SeqedTLVMsg{}
		err = conn.RecvMsg(msg)
		if err != nil {
			fmt.Println("read buf error:", err)
			return
		}

		fmt.Printf("read: %d %s\n", msg.Serial(), string(msg.Body()))

		time.Sleep(time.Duration(1<<id) * time.Second)
	}
}
