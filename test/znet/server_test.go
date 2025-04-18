package znet_test

import (
	"fmt"
	"my-zinx/zinx/znet"
	"net"
	"testing"
	"time"
)

func ClientTest() {
	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

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

		fmt.Printf(" server call back: %s\n", buf)

		time.Sleep(1 * time.Second)
	}
}

func TestServer(t *testing.T) {
	s := znet.NewServer("[zinx V0.2]", 3333) // 创建一个 server Handler

	go ClientTest() // 启动客户端

	s.Start()
	s.Serve()
}
