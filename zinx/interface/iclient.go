package iface

import (
	"context"
	"net"
)

type IClient interface {
	Connect() error
	Start(parent context.Context, fns ...func())
	Close()
	Conn() net.TCPConn
	SendMsg(msg IPacket) error
	RecvMsg(msg IPacket) error
}
