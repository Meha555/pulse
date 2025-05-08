package task

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	AddJobTag = iota
	SubJobTag
	MulJobTag
	DivJobTag
)

var KindStr = [...]byte{
	'+', '-', '*', '/',
}

// 本例采用二进制协议而非文本协议

type Request struct {
	ID uuid.UUID
	A  uint32
	B  uint32
}

// MarshalBinary implements encoding.BinaryMarshaler
func (r Request) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	// 复制一份内存
	uuid, _ := r.ID.MarshalBinary()
	if err := binary.Write(buf, binary.NativeEndian, uuid); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.NativeEndian, r.A); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.NativeEndian, r.B); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (r *Request) UnmarshalBinary(data []byte) error {
	if len(data) < 24 {
		return errors.New("data length is less than Request size")
	}
	if err := r.ID.UnmarshalBinary(data[:16]); err != nil {
		return fmt.Errorf("parse ID error: %w", err)
	}
	r.A = binary.NativeEndian.Uint32(data[16:20])
	r.B = binary.NativeEndian.Uint32(data[20:24])
	return nil
}

type Response struct {
	ID  uuid.UUID
	Res uint32
}

func (r *Response) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	uuid, _ := r.ID.MarshalBinary()
	if err := binary.Write(buf, binary.NativeEndian, uuid); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.NativeEndian, r.Res); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *Response) UnmarshalBinary(data []byte) error {
	if len(data) < 20 {
		return errors.New("data length is less than Response size")
	}
	if err := r.ID.UnmarshalBinary(data[:16]); err != nil {
		return fmt.Errorf("parse ID error: %w", err)
	}
	r.Res = binary.NativeEndian.Uint32(data[16:20])
	return nil
}
