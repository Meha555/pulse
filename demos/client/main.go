package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Client Test ... start")

	conn, err := net.Dial("tcp4", "127.0.0.1:3333")
	if err != nil {
		fmt.Println("client start err: ", err)
		return
	}

	for {
		cnt, err := conn.Write([]byte("hello ZINX"))
		if err != nil {
			fmt.Println("Write error err", err)
			return
		}

		buf := make([]byte, cnt)
		_, err = conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		fmt.Printf("read: %s\n", buf)

		time.Sleep(1 * time.Second)
	}
}
