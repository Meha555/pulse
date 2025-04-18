package znet_test

import (
	"my-zinx/zinx/znet"
	"testing"
)

func TestPacket(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		p := znet.NewPacket([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		data, err := znet.Marshal(p)
		if err != nil {
			t.Fatalf("Marshal Packet 失败: %v", err)
		}
		if len(data) != 14 {
			t.Errorf("Marshal 结果长度不正确，期望 14，实际 %d", len(data))
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		data := []byte{10, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		p := &znet.Packet{}
		err := znet.Unmarshal(data, p, true)
		if err != nil {
			t.Fatalf("Unmarshal Packet 失败: %v", err)
		}
		if p.BodyLen() != 10 {
			t.Errorf("bodyLen 不正确，期望 10，实际 %d", p.BodyLen())
		}
		if string(p.Body()) != string([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
			t.Errorf("body 不正确")
		}
	})
}

func TestTLVMsg(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		tlv := znet.NewTLVMsg(1, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		data, err := znet.Marshal(tlv)
		if err != nil {
			t.Fatalf("Marshal TLVMsg 失败: %v", err)
		}
		if len(data) != 16 {
			t.Errorf("Marshal 结果长度不正确，期望 18，实际 %d", len(data))
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		data := []byte{1, 0, 10, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		tlv := &znet.TLVMsg{}
		err := znet.Unmarshal(data, tlv, true)
		if err != nil {
			t.Fatalf("Unmarshal TLVMsg 失败: %v", err)
		}
		if tlv.Tag() != 1 {
			t.Errorf("tag 不正确，期望 1，实际 %d", tlv.Tag())
		}
		if tlv.BodyLen() != 10 {
			t.Errorf("bodyLen 不正确，期望 10，实际 %d", tlv.BodyLen())
		}
		if string(tlv.Body()) != string([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
			t.Errorf("body 不正确")
		}
	})
}

func TestSeqedMsg(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		seq := znet.NewSeqedMsg(1, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		data, err := znet.Marshal(seq)
		if err != nil {
			t.Fatalf("Marshal SeqedMsg 失败: %v", err)
		}
		if len(data) != 18 {
			t.Errorf("Marshal 结果长度不正确，期望 18，实际 %d", len(data))
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		data := []byte{1, 0, 0, 0, 10, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		seq := &znet.SeqedMsg{}
		err := znet.Unmarshal(data, seq, true)
		if err != nil {
			t.Fatalf("Unmarshal SeqedMsg 失败: %v", err)
		}
		if seq.Serial() != 1 {
			t.Errorf("serial 不正确，期望 1，实际 %d", seq.Serial())
		}
		if seq.BodyLen() != 10 {
			t.Errorf("bodyLen 不正确，期望 10，实际 %d", seq.BodyLen())
		}
		if string(seq.Body()) != string([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
			t.Errorf("body 不正确")
		}
	})
}
