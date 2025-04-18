package ziface

import (
	"net"

	"github.com/google/uuid"
)

// IConnection
// 与具体业务回调绑定的TCP连接
type IConnection interface {
	// 让当前连接开始工作
	Open()
	// 停止连接的工作，关闭连接
	Close()
	// 获取该对象的唯一标识
	ConnID() uuid.UUID
	RemoteAddr() net.Addr
}

// 所有Connection在处理业务时的钩子方法的函数签名
type HandleFunc func(peer *net.TCPConn, data []byte, cnt int) error
