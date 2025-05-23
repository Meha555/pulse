package session

import (
	"context"
	"my-zinx/core/message"
	"my-zinx/server/common"
)

type Request struct {
	// 已经和客户端建立好的连接会话
	session common.ISession
	// 客户端请求的数据（要求是顺序的TLV消息）
	msg message.ISeqedTLVMsg
	// 传递的参数（上下文）
	valueCtx context.Context
}

func NewRequest(conn common.ISession, msg message.ISeqedTLVMsg) *Request {
	return &Request{
		session:  conn,
		msg:      msg,
		valueCtx: context.Background(),
	}
}

func (r *Request) Session() common.ISession {
	return r.session
}

func (r *Request) Msg() message.ISeqedTLVMsg {
	return r.msg
}

func (r *Request) Set(key string, value interface{}) {
	r.valueCtx = context.WithValue(r.valueCtx, key, value)
}

func (r *Request) Get(key string) (value interface{}, exists bool) {
	value = r.valueCtx.Value(key)
	return value, value != nil
}
