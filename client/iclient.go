package client

import (
	"context"
	"net"

	"github.com/Meha555/pulse/core/message"
)

type IClient interface {
	Connect() error
	Start(parent context.Context, fns ...func())
	Close()
	Conn() net.TCPConn
	SendMsg(msg message.IPacket) error
	RecvMsg(msg message.IPacket) error
}
