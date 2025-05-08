package iface

import (
	"net"

	"github.com/google/uuid"
)

// ISession
// 与具体业务回调绑定的TCP连接
type ISession interface {
	// 让当前连接开始工作
	Open() error
	// 停止连接的工作，关闭连接
	Close()
	// 获取该对象的唯一标识
	SessionID() uuid.UUID
	// 获取底层的socket
	Conn() net.Conn
	// 更新心跳次数
	UpdateHeartBeat()
	// 获取心跳计次
	HeartBeat() uint
	// 获取可读写的退出chan
	ExitChan() <-chan struct{}

	SendMsg(msg IPacket) error
	RecvMsg(msg IPacket) error
}

// 所有Connection在处理业务时的钩子方法的函数签名
// type HandleFunc func(peer *net.TCPConn, data []byte, cnt int) error
