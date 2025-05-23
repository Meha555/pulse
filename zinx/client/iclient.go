package client

import (
	"context"
	"my-zinx/core/message"
	"net"
)

type IClient interface {
	Connect() error
	Start(parent context.Context, fns ...func())
	Close()
	Conn() net.TCPConn
	SendMsg(msg message.IPacket) error
	RecvMsg(msg message.IPacket) error
}
