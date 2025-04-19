package znet_test

import (
	"bytes"
	"log"
	"regexp"
	"testing"

	"my-zinx/zinx/ziface"
	"my-zinx/zinx/znet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMsgQueue 模拟消息队列
type MockMsgQueue struct {
	mock.Mock
}

func (m *MockMsgQueue) Push(request ziface.IRequest) {
	m.Called(request)
}

func (m *MockMsgQueue) Pop() ziface.IRequest {
	args := m.Called()
	return args.Get(0).(ziface.IRequest)
}

func (m *MockMsgQueue) Len() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockMsgQueue) Cap() int {
	args := m.Called()
	return args.Int(0)
}

// MockApiMapper 模拟 API 映射器
type MockApiMapper struct {
	mock.Mock
}

func (m *MockApiMapper) ExecJob(tag uint16, request ziface.IRequest) error {
	args := m.Called(tag, request)
	return args.Error(0)
}

func (m *MockApiMapper) AddJob(tag uint16, job ziface.IJob) ziface.IApiMapper {
	m.Called(tag, job)
	return m
}

func (m *MockApiMapper) GetJob(tag uint16) ziface.IJob {
	args := m.Called(tag)
	return args.Get(0).(ziface.IJob)
}

// MockIRequest 模拟 IRequest 接口
type MockIRequest struct {
	mock.Mock
}

func (m *MockIRequest) Msg() ziface.ISeqedTLVMsg {
	args := m.Called()
	return args.Get(0).(ziface.ISeqedTLVMsg)
}

func (m *MockIRequest) Conn() ziface.IConnection {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(ziface.IConnection)
}

// MockIMessage 模拟 IMessage 接口
type MockIMessage struct {
	mock.Mock
}

func (m *MockIMessage) Tag() string {
	args := m.Called()
	return args.String(0)
}

func TestWorkerPool(t *testing.T) {
	// Redirect log output to buffer for verification
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(nil)

	// Create mocks
	mockMQ := new(MockMsgQueue)
	mockMapper := new(MockApiMapper)
	mockRequest := new(MockIRequest)
	mockMessage := new(MockIMessage)

	// Set up mock behaviors
	mockRequest.On("Msg").Return(mockMessage)
	mockRequest.On("Conn").Return(nil)
	mockMessage.On("Tag").Return("test-tag")
	mockMapper.On("ExecJob", uint16(0), mockRequest).Return(nil)

	// Initialize WorkerPool
	pool := znet.NewWokerPool(2, mockMQ, mockMapper)

	// Test Start and Stop
	t.Run("Start and Stop", func(t *testing.T) {
		mockMQ.On("Pop").Return(mockRequest).Once()
		mockMQ.On("Pop").Return(nil) // Simulate no more messages

		pool.Start()
		pool.Stop()

		// Verify logs using regex
		assert.Regexp(t, regexp.MustCompile(`Worker\[\d+\] started`), logBuffer.String())
		assert.Contains(t, logBuffer.String(), "All workers stopped")
	})

	// Test Post
	t.Run("Post", func(t *testing.T) {
		mockMQ.On("Push", mockRequest).Return()

		pool.Post(mockRequest)

		// Check if methods are called at least once
		mockMQ.AssertNumberOfCalls(t, "Push", 1)
	})
}