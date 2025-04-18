package znet

import "my-zinx/zinx/ziface"

type Request struct {
	// 已经和客户端建立好的连接
	conn ziface.IConnection
	// 客户端请求的数据
	data []byte
}

func NewRequest(conn ziface.IConnection) ziface.IRequest {
	return &Request{
		conn: conn,
		data: nil,
	}
}

func (r *Request) Conn() ziface.IConnection {
	return r.conn
}

func (r *Request) Data() []byte {
	return r.data
}

// 确保 Request 实现了 ziface.IRequest 接口
var _ ziface.IRequest = (*Request)(nil)
