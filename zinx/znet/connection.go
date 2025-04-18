package znet

import (
	"log"
	"my-zinx/zinx/ziface"
	"net"

	"github.com/google/uuid"
)

type Connection struct {
	// 当前连接的socket TCP套接字
	conn *net.TCPConn
	// 当前连接的ID 也可以称作为SessionID，ID全局唯一
	connID uuid.UUID
	// 当前连接的关闭状态
	isClosed bool
	// 回调函数
	handler ziface.HandleFunc
	// FIXME: 通知该连接已经停止 这个字段是否需要？
	// 为什么不直接在 Stop() 方法中调用 Conn.Close() 来关闭连接？
	exitChan chan struct{}
}

func NewConnection(conn *net.TCPConn, connID uuid.UUID, handler ziface.HandleFunc) *Connection {
	return &Connection{
		conn:     conn,
		connID:   connID,
		isClosed: false,
		handler:  handler,
		exitChan: make(chan struct{}, 1), // 这里设置为 1，确保至少有一个缓冲区，防止写入时没人读导致阻塞，或者反之
	}
}

func (c *Connection) Open() {
	go c.StartReader()

	for {
		select {
		case <-c.exitChan: // 等待 Stop() 方法通知退出
			return
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
		buf := make([]byte, 512)
		nbytes, err := c.conn.Read(buf)
		if err != nil {
			log.Println("Read error:", err)
			c.exitChan <- struct{}{} // 通知 Open() 方法退出
			continue                 // 允许这次没拿到数据，下次再拿
		}
		// 调用handler指定的业务处理数据
		if err := c.handler(c.conn, buf, nbytes); err != nil {
			log.Println("Handle error:", err)
			c.exitChan <- struct{}{} // 通知 Open() 方法退出
			return
		}
	}
}

// StartWriter 用于向客户端发送数据的 Goroutine
func (c *Connection) StartWriter() {

}
