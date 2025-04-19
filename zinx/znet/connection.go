package znet

import (
	"context"
	"errors"
	"fmt"
	"io"
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
	// 用于控制连接的超时
	ctx context.Context
	// 当前连接的关闭状态
	isClosed bool

	// 当前连接的回调函数
	// handler ziface.HandleFunc

	// 路由对象，用于处理当前连接的业务逻辑
	mapper ziface.IControllerMapper
	// 用于读写协程(Reader/Writer)之间的通信（用于实现读写业务分离）
	msgChan chan []byte
	// FIXME: 通知该连接已经停止（Reader通知Writer，因为对端关闭连接后Reader会收到EOF REVIEW 底层收到FIN，上报EOF）
	// 为什么不直接在 Stop() 方法中调用 Conn.Close() 来关闭连接？
	exitChan chan struct{}
}

func NewConnection(conn *net.TCPConn, mapper ziface.IControllerMapper) *Connection {
	return &Connection{
		conn:     conn,
		connID:   uuid.New(),
		ctx:      context.Background(),
		isClosed: false,
		mapper:   mapper,
		msgChan:  make(chan []byte, 10),  // 这里设置缓冲区大小为10，允许读写协程的处理速率有一定的差异
		exitChan: make(chan struct{}, 1), // 这里设置为 1，确保至少有一个缓冲区，防止写入时没人读导致阻塞，或者反之
	}
}

func (c *Connection) Open() error {
	if c.mapper == nil {
		return errors.New("controller is nil")
	}
	// 启动IO协程负责该连接的读写操作
	go c.Reader()
	go c.Writer()

	// 等待 Stop() 方法通知退出
	for range c.exitChan {
		return nil
	}
	return nil
}

func (c *Connection) Close() {
	if c.isClosed {
		return
	}
	c.isClosed = true

	// TODO 如果用户注册了该连接的关闭回调业务, 那么应该在此刻显式调用

	c.conn.Close()
	c.exitChan <- struct{}{} // 通知 Open() 方法退出
	close(c.msgChan)
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

func (c Connection) SendMsg(msg ziface.IPacket) error {
	if c.isClosed {
		return errors.New("connection is closed")
	}
	data, err := Marshal(msg)
	if err != nil {
		return err
	}
	// return c.conn.Write(data)
	// 提交给让Writer协程异步发送，这样不会因为底层TCP发送缓冲区满而导致这里阻塞
	// TODO 如果发送有错误，则由Writer协程处理
	c.msgChan <- data
	return nil
}

// TODO 这种接口作为传出参数，不用指针能否实现传出修改？
func (c Connection) RecvMsg(msg ziface.IPacket) error {
	if c.isClosed {
		return errors.New("connection is closed")
	}
	headerData := make([]byte, msg.HeaderLen())
	if _, err := io.ReadFull(c.conn, headerData); err != nil {
		return fmt.Errorf("read header error: %v", err)
	}
	if err := Unmarshal(headerData, msg, false); err != nil {
		return fmt.Errorf("unmarshal header err: %v", err)
	}
	// 读取负载
	if msg.BodyLen() <= 0 {
		return nil
	}
	bodyData := make([]byte, msg.BodyLen())
	if _, err := io.ReadFull(c.conn, bodyData); err != nil {
		return fmt.Errorf("read body error: %v", err)
	}
	if err := UmarshalBodyOnly(bodyData, int(msg.BodyLen()), msg); err != nil {
		return fmt.Errorf("Unmarshal body error: %v", err)
	}
	return nil
}

func (c Connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// 确保 Connection 实现 ziface.IConenction 方法
var _ ziface.IConnection = (*Connection)(nil)

// Reader 是用于读取客户端数据的 Goroutine
// 会需要与主协程通过chan通信
func (c *Connection) Reader() {
	log.Println("Reader Goroutine is running")
	defer log.Println(c.RemoteAddr().String(), " Reader Goroutine exit!")
	defer c.Close() // 确保连接能被关闭

	for {
		msg := &SeqedTLVMsg{}
		if err := c.RecvMsg(msg); err != nil {
			log.Println("RecvMsg error:", err)
			c.exitChan <- struct{}{} // 通知 Open() 方法退出
			return
		}
		// 封装请求数据
		req := NewRequest(c, msg)
		// 起协程来执行业务
		go c.mapper.ExecController(msg.Tag(), req)
	}
}

// Writer 是用于向客户端发送数据的 Goroutine
// 会需要与主协程通过chan通信
func (c *Connection) Writer() {
	log.Println("Writer Goroutine is running")
	defer log.Println(c.RemoteAddr().String(), " Writer Goroutine exit!")
	for {
		select {
		case data := <-c.msgChan: // 从msgChan 中读取数据
			if _, err := c.Send(data); err != nil {
				log.Println("Send error:", err)
				continue
			}
		case <-c.exitChan: // 响应退出信号
			return
		}
	}
}
