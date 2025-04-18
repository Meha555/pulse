package ziface

// Message表示Request中的数据。之前是用 []byte 来在 Request 中接收全部数据，这个结构过于简单。
// 如果我们能从 Request 结构当中得知消息的类型、长度，那就更好了。

// TLV(Tag-Len-Value/Body)消息
type ITLVMsg interface {
	IPacket
	// 消息的唯一Tag，固定2B
	Tag() uint16
	SetTag(tag uint16)
	// Body部分的长度（单位B）
	BodyLen() uint32
}

// 顺序消息
type ISequentialMsg interface {
	IPacket
	// 消息序列号（用于维护应用层消息的顺序）。固定4B
	Serial() uint32
	SetSerial(serial uint32)
}