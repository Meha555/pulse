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

	var serial uint32 = 0
	for {

		_, err := conn.SendMsg(znet.NewSeqedMsg(serial, []byte("hello ZINX")))
		if err != nil {
			fmt.Println("Write error err", err)
			return
		}
		serial++

		msg := &znet.SeqedMsg{}
		err = conn.RecvMsg(msg)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		fmt.Printf("read: %d %s\n", msg.Serial(), string(msg.Body()))

		time.Sleep(1 * time.Second)
	}
}
