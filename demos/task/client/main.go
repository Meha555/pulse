package main

import (
	"context"
	"errors"
	"example/demos/task"
	"math/rand/v2"
	"my-zinx/client"
	"my-zinx/core/job"
	"my-zinx/core/message"
	. "my-zinx/log"
	"my-zinx/utils"
	"time"

	"github.com/google/uuid"
)

var mq = utils.NewBlockingQueue[client.TaskHandler](5)
var pool = client.NewWorkerPool(3, mq)

func main() {
	pool.Start()
	cli := client.NewClient("127.0.0.1", 3333, client.WithExitTimeout(5), client.WithHeartBeatInterval(1))
	cli.Start(context.Background(), // NOTE 这里的多个func，实际上应该是作为线程池执行的任务，而不是直接作为一个线程
		func() { doReceive(cli) },
		func() { doCaculate(cli, task.AddJobTag) },
		func() { doCaculate(cli, task.SubJobTag) },
		func() { doCaculate(cli, task.MulJobTag) },
		func() { doCaculate(cli, task.DivJobTag) })
	pool.Stop()
	Log.Info("Client exit")
}

func doReceive(cli *client.Client) {
	for {
		msg := &message.SeqedTLVMsg{}
		if err := cli.RecvMsg(msg); err != nil {
			Log.Errorf("worker read buf error: %v", err)
			return
		}
		// 丢弃无关的心跳包
		if msg.Tag() == job.HeartBeatTag {
			continue
		}

		var rsp task.Response
		if err := rsp.UnmarshalBinary(msg.Body()); err != nil {
			Log.Errorf("worker unmarshal data error: %v", err)
			continue
		}

		task := client.GetTask(rsp.ID)
		if task == nil {
			Log.Errorf("task not found: %v", rsp.ID)
			continue
		}
		task.AppendData(rsp)
		task.Exec()
	}
}

func doCaculate(cli *client.Client, kind uint16) {
	var serial uint32 = 0
	for {
		taskID := uuid.New()
		A := rand.Uint32N(100)
		B := rand.Uint32N(100)
		if B > A {
			A, B = B, A
		}

		var (
			arg = task.Request{
				ID: taskID,
				A:  A,
				B:  B,
			}
			buf []byte
			err error
		)
		if buf, err = arg.MarshalBinary(); err != nil {
			Log.Errorf("write arg(%+v) failed: %v", arg, err)
			continue
		}
		msgSent := message.NewSeqedTLVMsg(serial, kind, buf)
		if err = cli.SendMsg(msgSent); err != nil {
			Log.Errorf("Write error: %v", err)
			return
		}
		Log.Infof("Send A = %d, B = %d, kind = %c", A, B, task.KindStr[kind])
		serial++

		client.NewTask(taskID, func(t *client.Task) error {
			if len(t.Data()) == 0 {
				return errors.New("task must have data")
			}
			A := t.Data()[0].(uint32)
			B := t.Data()[1].(uint32)
			rsp := t.Data()[2].(task.Response)
			Log.Infof("Task[%s] Res: %d %c %d = %d\n", rsp.ID, A, task.KindStr[kind], B, rsp.Res)
			return nil
		}, client.WithWorkerPool(pool), client.WithData(A, B))

		time.Sleep(time.Duration(3) * time.Second)
	}
}
