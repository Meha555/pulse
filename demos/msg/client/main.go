package main

import (
	"fmt"
	"my-zinx/zinx/znet"
	"net"
)

func main() {
	//客户端goroutine，负责模拟粘包的数据，然后进行发送
	conn, err := net.Dial("tcp", "127.0.0.1:3333")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	defer conn.Close()

	msg1 := znet.NewSeqedMsg(0, []byte{'h', 'e', 'l', 'l', 'o'})
	sendData1, err := znet.Marshal(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	msg2 := znet.NewSeqedMsg(1, []byte{'w', 'o', 'r', 'l', 'd', '!', '!'})
	sendData2, err := znet.Marshal(msg2)
	if err != nil {
		fmt.Println("client temp msg2 err:", err)
		return
	}

	//将sendData1，和 sendData2 拼接一起，组成粘包
	sendData1 = append(sendData1, sendData2...)

	//向服务器端写数据
	conn.Write(sendData1)
}
