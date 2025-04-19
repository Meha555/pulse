package ziface

// 很显然，我们不能把业务处理的方法绑死在type HandlerFunc func(*net.TCPConn, []byte, int) error这种格式中，我们需要定一些interface{}来让用户填写任意格式的连接处理业务方法。

// IRequest
// 作为Api的数据源，封装来自客户端的请求消息，此后需要什么数据都可以从该对象中获取
type IRequest interface {
	// 获取连接
	Conn() IConnection
	// 获取请求数据
	Msg() ISeqedTLVMsg
}
