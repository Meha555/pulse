package session

import (
	"log"
	"my-zinx/zinx/core/job"
	"my-zinx/zinx/core/message"
	iface "my-zinx/zinx/interface"
	"my-zinx/zinx/utils"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ConnMgr
// 支持在添加连接时自动监听其 exitChan，并在 exitCh 关闭时自动删除连接
type ConnMgr struct {
	conns map[uuid.UUID]iface.ISession

	// 用于心跳检查的定时器
	heartBeatTicker *time.Ticker
	mtx             sync.RWMutex
	wg              sync.WaitGroup
}

func NewConnMgr() *ConnMgr {
	c := &ConnMgr{
		conns:           make(map[uuid.UUID]iface.ISession),
		heartBeatTicker: time.NewTicker(time.Duration(utils.Conf.Server.HeartBeatTick) * time.Second),
	}

	// 心跳检查
	go func() {
		for range c.heartBeatTicker.C {
			c.mtx.Lock()
			// 时刻到，检查心跳情况
			for _, conn := range c.conns {
				go func(conn iface.ISession) {
					if conn.HeartBeat() < 5 {
						conn.(*Session).heartbeat++
						conn.(*Session).SendMsg(message.NewSeqedTLVMsg(0, job.HeartBeatTag, nil))
					} else {
						// 说明已经5 * utils.Conf.Server.HeartBeatTick秒未收到该客户端的心跳包，判定该客户端已经掉线
						log.Printf("Conn %s is timeout, maybe offline", conn.ConnID())
						c.Del(conn.ConnID())
					}
				}(conn)
			}
			c.mtx.Unlock()
		}
	}()

	return c
}

func (c *ConnMgr) Add(conn iface.ISession) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if _, exists := c.conns[conn.ConnID()]; exists {
		return
	}
	c.conns[conn.ConnID()] = conn

	// Start a goroutine to listen on the exitCh
	c.wg.Add(1)
	go func(connID uuid.UUID) {
		defer c.wg.Done()
		<-conn.ExitChan()
		c.Del(connID)
	}(conn.ConnID())
}

func (c *ConnMgr) Del(connID uuid.UUID) {
	c.mtx.Lock() // 死锁了
	defer c.mtx.Unlock()
	if conn, exists := c.conns[connID]; exists {
		conn.Close()
		delete(c.conns, connID)
	}
}

func (c *ConnMgr) Get(connID uuid.UUID) iface.ISession {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.conns[connID]
}

func (c *ConnMgr) Count() uint {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return uint(len(c.conns))
}

func (c *ConnMgr) Clear() {
	c.mtx.Lock()
	for connID := range c.conns {
		if conn, exists := c.conns[connID]; exists {
			conn.Close()
			delete(c.conns, connID)
		}
	}
	c.mtx.Unlock()
	// Wait for all goroutines to finish
	c.wg.Wait()
}

var _ iface.ISessionMgr = (*ConnMgr)(nil)
