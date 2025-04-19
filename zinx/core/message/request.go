package message

import iface "my-zinx/zinx/interface"

type Request struct {
	// 已经和客户端建立好的连接
	conn iface.IConnection
	// 客户端请求的数据（要求是顺序的TLV消息）
	msg iface.ISeqedTLVMsg
}

func NewRequest(conn iface.IConnection, msg iface.ISeqedTLVMsg) *Request {
	return &Request{
		conn: conn,
		msg:  msg,
	}
}

func (r *Request) Conn() iface.IConnection {
	return r.conn
}

func (r *Request) Msg() iface.ISeqedTLVMsg {
	return r.msg
}

// 确保 Request 实现了 iface.IRequest 接口
var _ iface.IRequest = (*Request)(nil)
