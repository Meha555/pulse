package iface

import "github.com/google/uuid"

// ISessionMgr 连接管理器接口
// 通过连接管理器，可以统一管理连接（新建、删除、获取）
// 典型应用就是广播、踢人
type ISessionMgr interface {
	// 添加一个连接
	Add(conn ISession)
	// 删除指定的连接
	Del(connID uuid.UUID)
	// 获取指定的连接
	Get(connID uuid.UUID) ISession
	// 当前连接个数
	Count() uint
	// 关闭所有连接并清空
	Clear()
}
