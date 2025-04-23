package session

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"my-zinx/zinx/core/job"
	"my-zinx/zinx/core/message"
	iface "my-zinx/zinx/interface"
	utils "my-zinx/zinx/utils"
	"net"
	"sync/atomic"

	"github.com/google/uuid"
)

type zHooks struct {
	OnOpen     SessionHook
	OnClose    SessionHook
	BeforeSend SessionHook
	BeforeRecv SessionHook
	AfterSend  SessionHook
	AfterRecv  SessionHook
}

type SessionHook func(iface.ISession)
type zHookOpt func(c *Session)

// 定义一个空函数
var noOp SessionHook = func(iface.ISession) {}

func OnOpen(f SessionHook) zHookOpt {
	return func(c *Session) {
		c.hookStub.OnOpen = f
	}
}

func OnClose(f SessionHook) zHookOpt {
	return func(c *Session) {
		c.hookStub.OnClose = f
	}
}

func BeforeSend(f SessionHook) zHookOpt {
	return func(c *Session) {
		c.hookStub.BeforeSend = f
	}
}

func BeforeRecv(f SessionHook) zHookOpt {
	return func(c *Session) {
		c.hookStub.BeforeRecv = f
	}
}

func AfterSend(f SessionHook) zHookOpt {
	return func(c *Session) {
		c.hookStub.AfterSend = f
	}
}

func AfterRecv(f SessionHook) zHookOpt {
	return func(c *Session) {
		c.hookStub.AfterRecv = f
	}
}

// Session
// 将裸的TCP socket包装，将具体的业务与连接绑定
type Session struct {
	// 当前连接的socket TCP套接字
	conn *net.TCPConn
	// 当前连接的ID 也可以称作为SessionID，ID全局唯一
	connID uuid.UUID
	// 当前连接的关闭状态
	isClosed atomic.Bool
	// 保活心跳次数
	heartbeat uint

	// 工作协程池
	wokerPool *job.WokerPool

	// 用于读写协程(Reader/Writer)之间的通信（用于实现读写业务分离）
	msgCh chan []byte
	// FIXME: 通知该连接已经停止（Reader通知Writer，因为对端关闭连接后Reader会收到EOF REVIEW 底层收到FIN，上报EOF）
	// 为什么不直接在 Stop() 方法中调用 Conn.Close() 来关闭连接？
	exitCh chan struct{}

	hookStub  zHooks
	valuedCtx utils.Context
}

func NewConnection(conn *net.TCPConn, parent context.Context, wokerPool *job.WokerPool, opts ...zHookOpt) *Session {
	c := &Session{
		conn:      conn,
		connID:    uuid.New(),
		isClosed:  atomic.Bool{},
		heartbeat: 0,
		wokerPool: wokerPool,
		msgCh:     make(chan []byte, utils.Conf.Server.MaxMsgQueueSize), // 这里设置缓冲区大小为10，允许读写协程的处理速率有一定的差异
		exitCh:    make(chan struct{}, 1),                               // 这里设置为 1，确保至少有一个缓冲区，防止写入时没人读导致阻塞，或者反之
		hookStub: zHooks{
			OnOpen:     noOp,
			OnClose:    noOp,
			BeforeSend: noOp,
			BeforeRecv: noOp,
			AfterSend:  noOp,
			AfterRecv:  noOp,
		},
		valuedCtx: utils.NewContext(parent),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Session) Open() error {
	// 启动IO协程负责该连接的读写操作
	go c.Reader()
	go c.Writer()

	c.hookStub.OnOpen(c)

	// 等待 Stop() 方法通知退出
	for range c.exitCh {
		return nil
	}
	return nil
}

func (c *Session) Close() {
	if c.isClosed.Load() {
		return
	}
	c.isClosed.Store(true)

	c.hookStub.OnClose(c)

	c.conn.Close()
	c.exitCh <- struct{}{} // 通知 Open() 方法退出
	close(c.msgCh)
	// TODO close管道，读端会收到一个零值，写端会收到一个错误？
	close(c.exitCh)
}

func (c *Session) ConnID() uuid.UUID {
	return c.connID
}

func (c *Session) Conn() net.Conn {
	return c.conn
}

func (c *Session) UpdateHeartBeat() {
	c.heartbeat = 0
}

func (c *Session) HeartBeat() uint {
	return c.heartbeat
}

func (c *Session) Send(data []byte) (int, error) {
	if c.isClosed.Load() {
		return 0, errors.New("connection is closed")
	}
	return c.conn.Write(data)
}

func (c *Session) Recv(data []byte) (int, error) {
	if c.isClosed.Load() {
		return 0, errors.New("connection is closed")
	}
	return c.conn.Read(data)
}

func (c *Session) SendMsg(msg iface.IPacket) error {
	if c.isClosed.Load() {
		return errors.New("connection is closed")
	}
	c.hookStub.BeforeSend(c)
	defer c.hookStub.AfterSend(c)
	data, err := message.Marshal(msg)
	if err != nil {
		return err
	}
	// return c.conn.Write(data)
	// 提交给让Writer协程异步发送，这样不会因为底层TCP发送缓冲区满而导致这里阻塞
	// TODO 如果发送有错误，则由Writer协程处理
	c.msgCh <- data
	return nil
}

// TODO 这种接口作为传出参数，不用指针能否实现传出修改？
func (c *Session) RecvMsg(msg iface.IPacket) error {
	if c.isClosed.Load() {
		return errors.New("connection is closed")
	}
	c.hookStub.BeforeRecv(c)
	defer c.hookStub.AfterRecv(c)
	headerData := make([]byte, msg.HeaderLen())
	if _, err := io.ReadFull(c.conn, headerData); err != nil {
		return fmt.Errorf("read header error: %v", err)
	}
	if err := message.Unmarshal(headerData, msg, false); err != nil {
		return fmt.Errorf("unmarshal header err: %v", err)
	}
	// log.Printf("msg bodylen=%d", msg.BodyLen())
	// 读取负载
	if msg.BodyLen() <= 0 {
		return nil
	}
	bodyData := make([]byte, msg.BodyLen())
	if _, err := io.ReadFull(c.conn, bodyData); err != nil {
		return fmt.Errorf("read body error: %v", err)
	}
	if err := message.UmarshalBodyOnly(bodyData, int(msg.BodyLen()), msg); err != nil {
		return fmt.Errorf("unmarshal body error: %v", err)
	}
	return nil
}

func (c *Session) ExitChan() chan struct{} {
	return c.exitCh
}

// 确保 Connection 实现 iface.IConenction 方法
var _ iface.ISession = (*Session)(nil)

// Reader 是用于读取客户端数据的 Goroutine
// 会需要与主协程通过chan通信
func (c *Session) Reader() {
	log.Println("Reader Goroutine is running")
	defer log.Println(c.Conn().RemoteAddr().String(), " Reader Goroutine exit!")
	defer c.Close() // 确保连接能被关闭

	for {
		msg := &message.SeqedTLVMsg{}
		if err := c.RecvMsg(msg); err != nil {
			log.Println("RecvMsg error:", err)
			c.Close()
			return
		}
		// 封装请求数据
		req := message.NewRequest(c, msg)
		// 提交给协程池来处理业务
		c.wokerPool.Post(req)
	}
}

// Writer 是用于向客户端发送数据的 Goroutine
// 会需要与主协程通过chan通信
func (c *Session) Writer() {
	log.Println("Writer Goroutine is running")
	defer log.Println(c.Conn().RemoteAddr().String(), " Writer Goroutine exit!")
	for {
		select {
		case data := <-c.msgCh: // 从msgChan 中读取数据
			if _, err := c.Send(data); err != nil {
				log.Println("Send error:", err)
				continue
			}
		case <-c.exitCh: // 响应退出信号
			return
		}
	}
}
