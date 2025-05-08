package message

import (
	"context"
	iface "my-zinx/zinx/interface"
)

type Request struct {
	// 已经和客户端建立好的连接会话
	session iface.ISession
	// 客户端请求的数据（要求是顺序的TLV消息）
	msg iface.ISeqedTLVMsg
	// 传递的参数（上下文）
	valueCtx context.Context
}

func NewRequest(conn iface.ISession, msg iface.ISeqedTLVMsg) *Request {
	return &Request{
		session: conn,
		msg:     msg,
	}
}

func (r *Request) Session() iface.ISession {
	return r.session
}

func (r *Request) Msg() iface.ISeqedTLVMsg {
	return r.msg
}

func (r *Request) Set(key string, value interface{}) {
	r.valueCtx = context.WithValue(r.valueCtx, key, value)
}

func (r *Request) Get(key string) (value interface{}, exists bool) {
	value = r.valueCtx.Value(key)
	return value, value != nil
}

// 确保 Request 实现了 iface.IRequest 接口
var _ iface.IRequest = (*Request)(nil)
