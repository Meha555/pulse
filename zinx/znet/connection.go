package znet

import (
	"errors"
	"log"
	"my-zinx/zinx/ziface"
	"net"

	"github.com/google/uuid"
)

// Connection
// 将裸的TCP socket包装，将具体的业务与连接绑定
type Connection struct {
	// 当前连接的socket TCP套接字
	conn *net.TCPConn
	// 当前连接的ID 也可以称作为SessionID，ID全局唯一
	connID uuid.UUID
	// 当前连接的关闭状态
	isClosed bool

	// 当前连接的回调函数
	// handler ziface.HandleFunc

	// 路由对象，用于处理当前连接的业务逻辑
	router ziface.IRouter
	// FIXME: 通知该连接已经停止 这个字段是否需要？
	// 为什么不直接在 Stop() 方法中调用 Conn.Close() 来关闭连接？
	exitChan chan struct{}
}

func NewConnection(conn *net.TCPConn, connID uuid.UUID, router ziface.IRouter) *Connection {
	return &Connection{
		conn:     conn,
		connID:   connID,
		isClosed: false,
		router:   router,
		exitChan: make(chan struct{}, 1), // 这里设置为 1，确保至少有一个缓冲区，防止写入时没人读导致阻塞，或者反之
	}
}

func (c *Connection) Open() error {
	if c.router == nil {
		return errors.New("router is nil")
	}
	go c.StartReader()

	for {
		select {
		case <-c.exitChan: // 等待 Stop() 方法通知退出
			return nil
		}
	}
}

func (c *Connection) Close() {
	if c.isClosed {
		return
	}
	c.isClosed = true

	// TODO 如果用户注册了该连接的关闭回调业务, 那么应该在此刻显式调用
	c.conn.Close()
	c.exitChan <- struct{}{} // 通知 Open() 方法退出
	// TODO close管道，读端会收到一个零值，写端会收到一个错误？
	close(c.exitChan)
}

func (c Connection) ConnID() uuid.UUID {
	return c.connID
}

func (c Connection) Conn() net.Conn {
	return c.conn
}

func (c Connection) Send(data []byte) (int, error) {
	if c.isClosed {
		return 0, errors.New("connection is closed")
	}
	return c.conn.Write(data)
}

func (c Connection) Recv(data []byte) (int, error) {
	if c.isClosed {
		return 0, errors.New("connection is closed")
	}
	return c.conn.Read(data)
}

func (c Connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// 确保 Connection 实现 ziface.IConenction 方法
var _ ziface.IConnection = (*Connection)(nil)

// StartReader 用于读取客户端数据的 Goroutine
// 需要与主协程通过chan通信
func (c *Connection) StartReader() {
	log.Println("Reader Goroutine is running")
	defer log.Println(c.RemoteAddr().String(), " Reader Goroutine exit!")
	defer c.Close() // 确保连接能被关闭

	for {
		// 读取对端请求数据
		buf := make([]byte, 512)
		nbytes, err := c.conn.Read(buf)
		if err != nil {
			log.Println("Read error:", err)
			c.exitChan <- struct{}{} // 通知 Open() 方法退出
			continue                 // 允许这次没拿到数据，下次再拿
		}

		// 封装请求数据
		req := &Request{
			conn: c,
			data: buf[:nbytes],
		}
		// 执行指定的业务处理数据
		go func(request ziface.IRequest) {
			var err error
			if err = c.router.PreHandle(request); err != nil {
				log.Println("PreHandle error:", err)
				c.exitChan <- struct{}{} // 通知 Open() 方法退出
				return
			}
			if err = c.router.Handle(request); err != nil {
				log.Println("Handle error:", err)
				c.exitChan <- struct{}{} // 通知 Open() 方法退出
				return
			}
			if err = c.router.PostHandle(request); err != nil {
				log.Println("PostHandle error:", err)
				c.exitChan <- struct{}{}
				return
			}
		}(req)
	}
}

// StartWriter 用于向客户端发送数据的 Goroutine
func (c *Connection) StartWriter() {

}
