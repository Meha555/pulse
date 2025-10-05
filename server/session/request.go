package session

import (
	"context"
	"sync"

	"github.com/Meha555/pulse/core/message"
	"github.com/Meha555/pulse/server/common"
	"github.com/Meha555/pulse/utils"
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

var RequestPool *sync.Pool

func GetRequest(conn common.ISession, msg message.ISeqedTLVMsg) (req *Request) {
	if RequestPool == nil {
		req = NewRequest(conn, msg)
	} else {
		// 从池中去除的Request对象可能是已存在的，也可能是新构造的，所以都需要初始化
		req = RequestPool.Get().(*Request)
		req.session = conn
		req.msg = msg
		req.valueCtx = context.Background()
	}
	return
}

func PutRequest(request common.IRequest) {
	if RequestPool != nil {
		RequestPool.Put(request)
	}
}

func allocRequest() *Request {
	return &Request{
		session:  nil,
		msg:      nil,
		valueCtx: context.Background(),
	}
}

func init() {
	if utils.Conf.Server.RequestPoolMode {
		RequestPool = &sync.Pool{
			New: func() any {
				return allocRequest()
			},
		}
	}
}
