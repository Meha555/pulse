package main

import (
	"fmt"
	"io"
	"net"
	"pulse/core/message"
)

// 创建服务器gotoutine，负责从客户端goroutine读取粘包的数据，然后进行解析
func main() {
	//创建socket TCP Server
	listener, err := net.Listen("tcp", "127.0.0.1:3333")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server accept err:", err)
		}

		//启动协程处理客户端请求
		go func(conn net.Conn) {
			for {
				msg := &message.SeqedMsg{}
				//1 先读出流中的head部分
				headerData := make([]byte, msg.HeaderLen())
				_, err := io.ReadFull(conn, headerData) //ReadFull会读取正好len(headerData)个字节
				if err != nil {
					fmt.Println("read head error:", err)
					break
				}
				//将headData字节流 拆包到msg中
				err = message.Unmarshal(headerData, msg, false)
				if err != nil {
					fmt.Println("unmarshal header err:", err)
					continue
				}

				//msg 是有data数据的，需要再次读取data数据
				if msg.BodyLen() > 0 {
					bodyData := make([]byte, msg.BodyLen())
					_, err = io.ReadFull(conn, bodyData)
					if err != nil {
						fmt.Println("read body error")
						continue
					}
					err = message.UmarshalBodyOnly(bodyData, int(msg.BodyLen()), msg)
					if err != nil {
						fmt.Println("unmarshal body err:", err)
						continue
					}
					fmt.Printf("Recv msg: %d, %s\n", msg.Serial(), string(msg.Body()))
				}
			}
		}(conn)
	}

}
