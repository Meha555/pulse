package message

import (
	"bytes"
	"testing"
)

func TestPacket(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		p := NewPacket([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		data, err := Marshal(p)
		if err != nil {
			t.Fatalf("Marshal Packet 失败: %v", err)
		}
		if len(data) != 14 {
			t.Errorf("Marshal 结果长度不正确，期望 14，实际 %d", len(data))
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		data := []byte{10, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		p := &Packet{}
		err := Unmarshal(data, p, true)
		if err != nil {
			t.Fatalf("Unmarshal Packet 失败: %v", err)
		}
		if p.BodyLen() != 10 {
			t.Errorf("bodyLen 不正确，期望 10，实际 %d", p.BodyLen())
		}
		if !bytes.Equal(p.Body(), []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
			t.Errorf("body 不正确")
		}
	})
}

func TestTLVMsg(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		tlv := NewTLVMsg(1, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		data, err := Marshal(tlv)
		if err != nil {
			t.Fatalf("Marshal TLVMsg 失败: %v", err)
		}
		if len(data) != 16 {
			t.Errorf("Marshal 结果长度不正确，期望 18，实际 %d", len(data))
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		data := []byte{1, 0, 10, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		tlv := &TLVMsg{}
		err := Unmarshal(data, tlv, true)
		if err != nil {
			t.Fatalf("Unmarshal TLVMsg 失败: %v", err)
		}
		if tlv.Tag() != 1 {
			t.Errorf("tag 不正确，期望 1，实际 %d", tlv.Tag())
		}
		if tlv.BodyLen() != 10 {
			t.Errorf("bodyLen 不正确，期望 10，实际 %d", tlv.BodyLen())
		}
		if !bytes.Equal(tlv.Body(), []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
			t.Errorf("body 不正确")
		}
	})
}

func TestSeqedMsg(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		seq := NewSeqedMsg(1, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		data, err := Marshal(seq)
		if err != nil {
			t.Fatalf("Marshal SeqedMsg 失败: %v", err)
		}
		if len(data) != 18 {
			t.Errorf("Marshal 结果长度不正确，期望 18，实际 %d", len(data))
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		data := []byte{1, 0, 0, 0, 10, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		seq := &SeqedMsg{}
		err := Unmarshal(data, seq, true)
		if err != nil {
			t.Fatalf("Unmarshal SeqedMsg 失败: %v", err)
		}
		if seq.Serial() != 1 {
			t.Errorf("serial 不正确，期望 1，实际 %d", seq.Serial())
		}
		if seq.BodyLen() != 10 {
			t.Errorf("bodyLen 不正确，期望 10，实际 %d", seq.BodyLen())
		}
		if !bytes.Equal(seq.Body(), []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
			t.Errorf("body 不正确")
		}
	})
}

func TestSeqedTLVMsg(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		seqtlv := NewSeqedTLVMsg(1, 123, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		data, err := Marshal(seqtlv)
		if err != nil {
			t.Fatalf("Marshal SeqedTLVMsg 失败: %v", err)
		}
		if len(data) != 20 {
			t.Errorf("Marshal 结果长度不正确，期望 22，实际 %d", len(data))
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		data := []byte{1, 0, 0, 0, 123, 0, 10, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		seqtlv := &SeqedTLVMsg{}
		err := Unmarshal(data, seqtlv, true)
		if err != nil {
			t.Fatalf("Unmarshal SeqedTLVMsg 失败: %v", err)
		}
		if seqtlv.Serial() != 1 {
			t.Errorf("serial 不正确，期望 1，实际 %d", seqtlv.Serial())
		}
		if seqtlv.Tag() != 123 {
			t.Errorf("tag 不正确，期望 123，实际 %d", seqtlv.Tag())
		}
		if seqtlv.BodyLen() != 10 {
			t.Errorf("bodyLen 不正确，期望 10，实际 %d", seqtlv.BodyLen())
		}
		if !bytes.Equal(seqtlv.Body(), []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
			t.Errorf("body 不正确")
		}
	})
}
