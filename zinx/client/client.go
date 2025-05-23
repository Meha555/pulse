package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"my-zinx/core"
	"my-zinx/core/message"
	log "my-zinx/log"
	"my-zinx/server/job"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

var logger = log.NewStdLogger(log.LevelInfo, "client", "[%t] [%c %l] [%f:%C:%L:%g] %m", false)

type counter struct {
	count uint32
}

type Client struct {
	Name      string
	IPVersion string
	IP        string
	Port      uint16
	conn      *net.TCPConn

	heartBeatInterval time.Duration
	exitTimeout       time.Duration  // 超时时间，单位：秒
	wg                sync.WaitGroup // 用于等待所有协程退出，实现优雅退出

	serial counter
}

func NewClient(ip string, port uint16, opts ...ClientOptions) *Client {
	c := &Client{
		Name:      "MY-ZINX Client@" + uuid.New().String(),
		IPVersion: "tcp4",
		IP:        ip,
		Port:      port,
		conn:      nil,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Connect()

	return c
}

func (c *Client) Connect() error {
	if c.conn != nil {
		return errors.New("client conn is not nil, maybe already connected")
	}
	addr, err := net.ResolveTCPAddr(c.IPVersion, fmt.Sprintf("%s:%d", c.IP, c.Port))
	if err != nil {
		return err
	}
	conn, err := net.DialTCP(c.IPVersion, nil, addr)
	if err != nil {
		return err
	}
	c.conn = conn
	c.serial.count = 0
	logger.Infof("client connected to server %s:%d", c.IP, c.Port)
	return nil
}

// Start 启动客户端业务。
// 客户端业务由fn执行，所有由fn托管的业务逻辑可以保证客户端退出时业务已经结束（在不超时的情况下）
func (c *Client) Start(parent context.Context, fns ...func()) {
	logger.Info("client start with jobs...")
	// c.wg.Add(1)
	go c.heartBeat() // 不交由c.wg控制，因为心跳协程直接退出不会有影响
	for _, fn := range fns {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			fn()
		}()
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	sig := <-sigCh
	logger.Infof("Received signal: %v (%s), exiting in %v ...", sig, sig, c.exitTimeout*time.Second)

	ctx, cancel := context.WithTimeout(parent, c.exitTimeout*time.Second)
	defer cancel()

	quitCh := make(chan struct{})
	go func() {
		logger.Info("Wating for all goroutines to exit")
		c.wg.Wait()
		close(quitCh)
	}()

	select {
	case <-quitCh:
		logger.Info("All goroutines exited")
	case <-ctx.Done():
		logger.Info("Timeout, force exit")
	}
	c.Close()
}

func (c *Client) Close() {

	if c.conn != nil {
		c.conn.Close()
	}
	c.conn = nil
	c.serial.count = 0
}

func (c *Client) Conn() net.TCPConn {
	return *c.conn
}

func (c *Client) SendMsg(msg core.IPacket) error {
	if c.conn == nil {
		return errors.New("connection is closed")
	}
	data, err := message.Marshal(msg)
	if err != nil {
		return fmt.Errorf("client send msg marshal error: %w", err)
	}
	_, err = c.conn.Write(data)
	if err != nil {
		return fmt.Errorf("client send msg write error: %w", err)
	}
	c.serial.count++
	return nil
}

func (c *Client) RecvMsg(msg core.IPacket) error {
	if c.conn == nil {
		return errors.New("connection is closed")
	}
	headerData := make([]byte, msg.HeaderLen())
	if _, err := io.ReadFull(c.conn, headerData); err != nil {
		return fmt.Errorf("read header error: %w", err)
	}
	if err := message.Unmarshal(headerData, msg, false); err != nil {
		return fmt.Errorf("unmarshal header err: %w", err)
	}
	// 读取负载
	if msg.BodyLen() <= 0 {
		return nil
	}
	bodyData := make([]byte, msg.BodyLen())
	if _, err := io.ReadFull(c.conn, bodyData); err != nil {
		return fmt.Errorf("read body error: %w", err)
	}
	if err := message.UmarshalBodyOnly(bodyData, int(msg.BodyLen()), msg); err != nil {
		return fmt.Errorf("unmarshal body error: %w", err)
	}
	return nil
}

// HeartBeat 方法用于向服务器发送心跳消息。
// 参数 interval 表示心跳消息的发送间隔，单位为秒。
func (c *Client) heartBeat() {
	// defer c.wg.Done()
	ticker := time.NewTicker(c.heartBeatInterval * time.Second)
	for range ticker.C {
		msgSent := message.NewSeqedTLVMsg(c.serial.count, job.HeartBeatTag, nil)
		if err := c.SendMsg(msgSent); err != nil {
			logger.Errorf("Write error: %v", err)
			return
		}
	}
}
