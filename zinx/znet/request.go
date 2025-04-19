package znet

import "my-zinx/zinx/ziface"

type Request struct {
	// 已经和客户端建立好的连接
	conn ziface.IConnection
	// 客户端请求的数据（要求是顺序的TLV消息）
	msg ziface.ISeqedTLVMsg
}

func NewRequest(conn ziface.IConnection, msg ziface.ISeqedTLVMsg) *Request {
	return &Request{
		conn: conn,
		msg:  msg,
	}
}

func (r *Request) Conn() ziface.IConnection {
	return r.conn
}

func (r *Request) Msg() ziface.ISeqedTLVMsg {
	return r.msg
}

// 确保 Request 实现了 ziface.IRequest 接口
var _ ziface.IRequest = (*Request)(nil)
