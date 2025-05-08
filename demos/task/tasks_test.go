package task

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/google/uuid"
)

func TestArg_MarshalBinary(t *testing.T) {
	// 创建一个 Arg 实例
	arg := Request{
		ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		A:  100,
		B:  200,
	}

	// 序列化
	data, err := arg.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}

	// 验证序列化结果
	expected := make([]byte, 24)
	copy(expected[:16], arg.ID[:])
	binary.NativeEndian.PutUint32(expected[16:], uint32(arg.A))
	binary.NativeEndian.PutUint32(expected[20:], uint32(arg.B))

	if !bytes.Equal(data, expected) {
		t.Errorf("MarshalBinary output mismatch: got %v, want %v", data, expected)
	}
}

func TestArg_UnmarshalBinary(t *testing.T) {
	// 构造测试数据
	id := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	data := make([]byte, 24)
	copy(data[:16], id[:])
	binary.NativeEndian.PutUint32(data[16:], 100)
	binary.NativeEndian.PutUint32(data[20:], 200)

	// 反序列化
	var arg Request
	err := arg.UnmarshalBinary(data)
	if err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	// 验证反序列化结果
	if arg.ID != id || arg.A != 100 || arg.B != 200 {
		t.Errorf("UnmarshalBinary output mismatch: got %+v, want %+v", arg, Request{ID: id, A: 100, B: 200})
	}
}

func TestArg_UnmarshalBinary_Error(t *testing.T) {
	// 测试数据长度不足的情况
	data := make([]byte, 15) // 数据长度小于 24 字节
	var arg Request
	err := arg.UnmarshalBinary(data)
	if err == nil {
		t.Error("UnmarshalBinary should return an error for insufficient data length")
	}

	// 测试无效 UUID 的情况
	data = []byte{}
	copy(data, []byte("invalid_uuid"))
	err = arg.UnmarshalBinary(data)
	if err == nil {
		t.Error("UnmarshalBinary should return an error for invalid UUID")
	}
}
