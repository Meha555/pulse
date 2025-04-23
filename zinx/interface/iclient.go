package iface

type IClient interface {
	Connect() error
	Exit()
	Conn() ISession
	SendMsg(msg IPacket) error
	RecvMsg(msg IPacket) error
}
