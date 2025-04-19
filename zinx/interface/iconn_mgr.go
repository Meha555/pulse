package iface

import "github.com/google/uuid"

// IConnMgr 连接管理器接口
// 通过连接管理器，可以统一管理连接（新建、删除、获取）
// 典型应用就是广播、踢人
// TODO 还要支持自定义的连接管理能力，比如心跳
type IConnMgr interface {
	// 添加一个连接
	Add(conn IConnection)
	// 删除指定的连接
	Del(connID uuid.UUID)
	// 获取指定的连接
	Get(connID uuid.UUID) IConnection
	// 当前连接个数
	Count() uint
	// 关闭所有连接并清空
	Clear()
}
