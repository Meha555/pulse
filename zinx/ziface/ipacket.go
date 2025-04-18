package ziface

// IPacket 数据包
type IPacket interface {
	// 消息（Raw Body，是这次从socket读出的全部内容）
	Body() []byte
	SetBody(data []byte)

	BodyLen() uint32
	HeaderLen() uint32
}
