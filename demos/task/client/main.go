package main

import (
	"context"
	"fmt"
	"my-zinx/client"
	"my-zinx/core/message"
	. "my-zinx/log"
	"time"
)

func main() {
	cli := client.NewClient("127.0.0.1", 3333, client.WithExitTimeout(5), client.WithHeartBeatInterval(1))
	cli.Start(context.Background(),
		func() { doEcho(cli, 0) },
		func() { doEcho(cli, 1) })
}

func doEcho(cli *client.Client, id uint16) {
	var serial uint32 = 0
	for {
		msgSent := message.NewSeqedTLVMsg(serial, id, fmt.Appendf(nil, "hello ZINX [%d]", id))
		if err := cli.SendMsg(msgSent); err != nil {
			Log.Errorf("Write error: %v", err)
			return
		}
		serial++

		msg := &message.SeqedTLVMsg{}
		if err := cli.RecvMsg(msg); err != nil {
			Log.Errorf("read buf error: %v", err)
			return
		}

		Log.Infof("[%d] read: %d [%d] %s\n", id, msg.Serial(), msg.Tag(), string(msg.Body()))
		if msg.Tag() != id {
			// 服务端只是返回属于该连接的消息，至于消息到底给到哪个协程，这个是客户端需要实现的
			Log.Warnf("want: %d, got: %d", id, msg.Tag())
		}

		time.Sleep(time.Duration(1<<id) * time.Second)
	}
}
