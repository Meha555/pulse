package session

import "github.com/Meha555/pulse/server/common"

type hooks struct {
	onOpen     hook
	onClose    hook
	beforeSend hook
	beforeRecv hook
	afterSend  hook
	afterRecv  hook
}

type hook func(common.ISession)
type hookOpt func(c *Session)

// 定义一个空函数
var noOp hook = func(common.ISession) {}

func OnOpen(f hook) hookOpt {
	return func(c *Session) {
		c.hookStub.onOpen = f
	}
}

func OnClose(f hook) hookOpt {
	return func(c *Session) {
		c.hookStub.onClose = f
	}
}

func BeforeSend(f hook) hookOpt {
	return func(c *Session) {
		c.hookStub.beforeSend = f
	}
}

func BeforeRecv(f hook) hookOpt {
	return func(c *Session) {
		c.hookStub.beforeRecv = f
	}
}

func AfterSend(f hook) hookOpt {
	return func(c *Session) {
		c.hookStub.afterSend = f
	}
}

func AfterRecv(f hook) hookOpt {
	return func(c *Session) {
		c.hookStub.afterRecv = f
	}
}
