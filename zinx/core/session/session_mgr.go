package session

import (
	"my-zinx/zinx/core/job"
	"my-zinx/zinx/core/message"
	iface "my-zinx/zinx/interface"
	"my-zinx/zinx/utils"
	"sync"
	"time"

	"github.com/google/uuid"
)

// SessionMgr
// 支持在添加连接时自动监听其 exitChan，并在 exitCh 关闭时自动删除连接
type SessionMgr struct {
	sessionMap map[uuid.UUID]iface.ISession

	// 用于心跳检查的定时器
	heartBeatTicker *time.Ticker
	mtx             sync.RWMutex
	wg              sync.WaitGroup
}

func NewSessionMgr() *SessionMgr {
	c := &SessionMgr{
		sessionMap:        make(map[uuid.UUID]iface.ISession),
		heartBeatTicker: time.NewTicker(time.Duration(utils.Conf.Server.HeartBeatTick) * time.Second),
	}

	// 心跳检查
	go func() {
		for range c.heartBeatTicker.C {
			c.mtx.Lock()
			// 时刻到，检查心跳情况
			for _, session := range c.sessionMap {
				go func(session iface.ISession) {
					if session.HeartBeat() < 5 {
						session.(*Session).heartbeat++
						session.(*Session).SendMsg(message.NewSeqedTLVMsg(0, job.HeartBeatTag, nil))
					} else {
						// 说明已经5 * utils.Conf.Server.HeartBeatTick秒未收到该客户端的心跳包，判定该客户端已经掉线
						logger.Warnf("Conn %s is timeout, maybe offline", session.SessionID())
						c.Del(session.SessionID())
					}
				}(session)
			}
			c.mtx.Unlock()
		}
	}()

	return c
}

func (c *SessionMgr) Add(session iface.ISession) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if _, exists := c.sessionMap[session.SessionID()]; exists {
		return
	}
	c.sessionMap[session.SessionID()] = session

	// Start a goroutine to listen on the exitCh
	c.wg.Add(1)
	go func(sessionID uuid.UUID) {
		defer c.wg.Done()
		<-session.ExitChan()
		c.Del(sessionID)
	}(session.SessionID())
}

func (c *SessionMgr) Del(sessionID uuid.UUID) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if session, exists := c.sessionMap[sessionID]; exists {
		session.Close()
		delete(c.sessionMap, sessionID)
	}
}

func (c *SessionMgr) Get(sessionID uuid.UUID) iface.ISession {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.sessionMap[sessionID]
}

func (c *SessionMgr) Count() uint {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return uint(len(c.sessionMap))
}

func (c *SessionMgr) Clear() {
	c.mtx.Lock()
	for sessionID := range c.sessionMap {
		if session, exists := c.sessionMap[sessionID]; exists {
			session.Close()
			delete(c.sessionMap, sessionID)
		}
	}
	c.mtx.Unlock()
	// Wait for all goroutines to finish
	c.wg.Wait()
}

var _ iface.ISessionMgr = (*SessionMgr)(nil)
