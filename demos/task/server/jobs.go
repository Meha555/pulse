package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"my-zinx/zinx/core/message"
	"my-zinx/zinx/core/session"
	iface "my-zinx/zinx/interface"
	"unsafe"
)

const (
	AddJobTag = iota
	SubJobTag
	MulJobTag
	DivJobTag
)

// 本例采用二进制协议而非文本协议

type Arg struct {
	A int32
	B int32
}

func ParseArg(data []byte) (A int32, B int32, err error) {
	if len(data) < int(unsafe.Sizeof(Arg{})) {
		err = errors.New("data length is less than Arg size")
		return
	}
	A = int32(binary.NativeEndian.Uint32(data[:4]))
	B = int32(binary.NativeEndian.Uint32(data[4:8]))
	return
}

type CalculateJob struct {
	calculator func(int32, int32) int32
}

func (j *CalculateJob) PreHandle(req iface.IRequest) error {
	A, B, err := ParseArg(req.Msg().Body())
	if err != nil {
		return fmt.Errorf("parse args error: %w", err)
	}
	req.Set("A", A)
	req.Set("B", B)
	return nil
}

func (j *CalculateJob) Handle(req iface.IRequest) error {
	var A, B int32
	if a, ok := req.Get("A"); !ok {
		return errors.New("request.Get(\"A\") failed")
	} else {
		A = a.(int32)
	}
	if b, ok := req.Get("B"); !ok {
		return errors.New("request.Get(\"B\") failed")
	} else {
		B = b.(int32)
	}
	res := j.calculator(A, B)
	req.Set("Result", res)
	return nil
}

func (j *CalculateJob) PostHandle(req iface.IRequest) error {
	var res int32
	if val, ok := req.Get("Result"); !ok {
		return errors.New("request.Get(\"Result\") failed")
	} else {
		res = val.(int32)
	}
	// 将 res 转换为字节切片，以便写入消息体
	resBytes := make([]byte, 4)
	binary.NativeEndian.PutUint32(resBytes, uint32(res))
	msg := message.NewSeqedTLVMsg(req.Msg().Serial()+1, AddJobTag, resBytes)
	if err := req.Session().(*session.Session).SendMsg(msg); err != nil {
		return fmt.Errorf("response error: %w", err)
	} else {
		return nil
	}
}

type AddJob struct {
	CalculateJob
}

type SubJob struct {
	CalculateJob
}

type MulJob struct {
	CalculateJob
}

type DivJob struct {
	CalculateJob
}

// Simple Factory
type CalculateJobFactory struct{}

func (f *CalculateJobFactory) CreateCalculator(tag uint16) iface.IJob {
	switch tag {
	case AddJobTag:
		return &AddJob{
			CalculateJob: CalculateJob{calculator: func(a int32, b int32) int32 { return a + b }},
		}
	case SubJobTag:
		return &SubJob{
			CalculateJob: CalculateJob{calculator: func(a int32, b int32) int32 { return a - b }},
		}
	case MulJobTag:
		return &MulJob{
			CalculateJob: CalculateJob{calculator: func(a int32, b int32) int32 { return a * b }},
		}
	case DivJobTag:
		return &DivJob{
			CalculateJob: CalculateJob{calculator: func(a int32, b int32) int32 {
				return a / b
			}},
		}
	default:
		return nil
	}
}
