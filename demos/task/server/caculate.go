package main

import (
	"errors"
	"example/demos/task"
	"fmt"
	"my-zinx/core/message"
	"my-zinx/server/common"

	. "my-zinx/log"

	"github.com/google/uuid"
)

type CalculateJob struct {
	kind       uint16
	calculator func(uint32, uint32) uint32
}

func (j *CalculateJob) PreHandle(req common.IRequest) error {
	var arg task.Request
	err := arg.UnmarshalBinary(req.Msg().Body())
	if err != nil {
		return fmt.Errorf("parse args error: %w", err)
	}
	req.Set("ID", arg.ID)
	req.Set("A", arg.A)
	req.Set("B", arg.B)
	return nil
}

func (j *CalculateJob) Handle(req common.IRequest) error {
	var A, B uint32
	if a, ok := req.Get("A"); !ok {
		return errors.New("request.Get(\"A\") failed")
	} else {
		A = a.(uint32)
	}
	if b, ok := req.Get("B"); !ok {
		return errors.New("request.Get(\"B\") failed")
	} else {
		B = b.(uint32)
	}
	res := j.calculator(A, B)
	Log.Debugf("Res: %d %c %d = %d\n", A, task.KindStr[j.kind], B, res)
	req.Set("Result", res)
	return nil
}

func (j *CalculateJob) PostHandle(req common.IRequest) error {
	var ID uuid.UUID
	if val, ok := req.Get("ID"); !ok {
		return errors.New("request.Get(\"ID\") failed")
	} else {
		ID = val.(uuid.UUID)
	}
	var Res uint32
	if val, ok := req.Get("Result"); !ok {
		return errors.New("request.Get(\"Result\") failed")
	} else {
		Res = val.(uint32)
	}
	arg := task.Response{
		ID:  ID,
		Res: Res,
	}
	Log.Warnf("Res: %+v", arg)
	data, err := arg.MarshalBinary()
	if err != nil {
		return fmt.Errorf("response marshal error: %w", err)
	}
	msg := message.NewSeqedTLVMsg(req.Msg().Serial()+1, j.kind, data)
	if err := req.Session().SendMsg(msg); err != nil {
		return fmt.Errorf("response send error: %w", err)
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

func (f *CalculateJobFactory) CreateCalculator(tag uint16) common.IJob {
	switch tag {
	case task.AddJobTag:
		return &AddJob{
			CalculateJob: CalculateJob{kind: tag, calculator: func(a uint32, b uint32) uint32 { return a + b }},
		}
	case task.SubJobTag:
		return &SubJob{
			CalculateJob: CalculateJob{kind: tag, calculator: func(a uint32, b uint32) uint32 { return a - b }},
		}
	case task.MulJobTag:
		return &MulJob{
			CalculateJob: CalculateJob{kind: tag, calculator: func(a uint32, b uint32) uint32 { return a * b }},
		}
	case task.DivJobTag:
		return &DivJob{
			CalculateJob: CalculateJob{kind: tag, calculator: func(a uint32, b uint32) uint32 {
				if b == 0 {
					return 0
				}
				return a / b
			}},
		}
	default:
		return nil
	}
}
