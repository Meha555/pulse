package znet

import "my-zinx/zinx/ziface"

type Request struct {
	// 已经和客户端建立好的连接
	conn ziface.IConnection
	// 客户端请求的数据
	msg ziface.IPacket
}

func NewRequest(conn ziface.IConnection) ziface.IRequest {
	return &Request{
		conn: conn,
		msg:  nil,
	}
}

func (r *Request) Conn() ziface.IConnection {
	return r.conn
}

func (r *Request) Msg() ziface.IPacket {
	return r.msg
}

// 确保 Request 实现了 ziface.IRequest 接口
var _ ziface.IRequest = (*Request)(nil)
