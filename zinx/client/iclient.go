package client

import (
	"context"
	"my-zinx/core"
	"net"
)

type IClient interface {
	Connect() error
	Start(parent context.Context, fns ...func())
	Close()
	Conn() net.TCPConn
	SendMsg(msg core.IPacket) error
	RecvMsg(msg core.IPacket) error
}
