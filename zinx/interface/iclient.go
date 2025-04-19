package iface

type IClient interface {
	Connect() error
	Exit()
	Conn() IConnection
	SendMsg(msg IPacket) error
	RecvMsg(msg IPacket) error
}
