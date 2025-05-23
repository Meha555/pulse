package message

// IPacket 数据包
type IPacket interface {
	// 消息（Raw Body，是这次从socket读出的全部内容）
	Body() []byte
	SetBody(data []byte)

	BodyLen() uint32
	HeaderLen() uint32
}

// Message表示Request中的数据。之前是用 []byte 来在 Request 中接收全部数据，这个结构过于简单。
// 如果我们能从 Request 结构当中得知消息的类型、长度，那就更好了。

// TLV(Tag-Len-Value/Body)消息
type ITLVMsg interface {
	IPacket
	// 消息的唯一Tag，固定2B
	Tag() uint16
	SetTag(tag uint16)
}

// 顺序消息
type ISeqedMsg interface {
	IPacket
	// 消息序列号（用于维护应用层消息的顺序）。固定4B
	Serial() uint32
	SetSerial(serial uint32)
}

type ISeqedTLVMsg interface {
	ITLVMsg
	ISeqedMsg
}

// Packet
// NOTE 应该被嵌入到结构体的最后
//
//		+------------+
//	 | Len | Body |
//	 +------------+
type Packet struct {
	bodyLen uint32
	body    []byte
}

func NewPacket(data []byte) *Packet {
	var packet Packet
	packet.SetBody(data)
	return &packet
}

func (p Packet) Body() []byte {
	return p.body
}

func (p *Packet) SetBody(body []byte) {
	p.bodyLen = uint32(len(body))
	p.body = body
}

func (p Packet) BodyLen() uint32 {
	return p.bodyLen
}

func (p Packet) HeaderLen() uint32 {
	return 4 // sizeof(uint32)
}

// TLVMsg
//
//		+------------------+
//	 | Tag | Len | Body |
//	 +------------------+
type TLVMsg struct {
	tag uint16
	Packet
}

func NewTLVMsg(tag uint16, data []byte) *TLVMsg {
	var msg TLVMsg
	msg.SetTag(tag)
	msg.SetBody(data)
	return &msg
}

func (t TLVMsg) Tag() uint16 {
	return t.tag
}

func (t *TLVMsg) SetTag(tag uint16) {
	t.tag = tag
}

func (t TLVMsg) HeaderLen() uint32 {
	return 2 + 4 // sizeof(uint16) + sizeof(uint32)
}

// SeqMsg顺序的TLV消息
// +---------------------+
// | Serial | Len | Body |
// +---------------------+
type SeqedMsg struct {
	serial uint32
	Packet
}

func NewSeqedMsg(serial uint32, data []byte) *SeqedMsg {
	var msg SeqedMsg
	msg.SetSerial(serial)
	msg.SetBody(data)
	return &msg
}

func (s SeqedMsg) Serial() uint32 {
	return s.serial
}

func (s *SeqedMsg) SetSerial(serial uint32) {
	s.serial = serial
}

func (s SeqedMsg) HeaderLen() uint32 {
	return 4 + 4 // sizeof(uint32) + sizeof(uint32)
}

type SeqedTLVMsg struct {
	serial uint32
	TLVMsg
}

func NewSeqedTLVMsg(serial uint32, tag uint16, data []byte) *SeqedTLVMsg {
	var msg SeqedTLVMsg
	msg.SetSerial(serial)
	msg.SetTag(tag)
	msg.SetBody(data)
	return &msg
}

func (s SeqedTLVMsg) Serial() uint32 {
	return s.serial
}

func (s *SeqedTLVMsg) SetSerial(serial uint32) {
	s.serial = serial
}

func (s SeqedTLVMsg) HeaderLen() uint32 {
	return 4 /*sizeof(uint32)*/ + s.TLVMsg.HeaderLen()
}

// TODO 利用io.Reader和io.Writer接口、Proxy设计模式封装读写msg的操作
