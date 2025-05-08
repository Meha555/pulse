package core_test

import (
	"bytes"
	"log"
	"regexp"
	"testing"

	"my-zinx/core/job"
	iface "my-zinx/interface"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMsgQueue 模拟消息队列
type MockMsgQueue struct {
	mock.Mock
}

func (m *MockMsgQueue) Push(request iface.IRequest) {
	m.Called(request)
}

func (m *MockMsgQueue) Pop() iface.IRequest {
	args := m.Called()
	return args.Get(0).(iface.IRequest)
}

func (m *MockMsgQueue) Len() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockMsgQueue) Cap() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockMsgQueue) Close() {
	m.Called()
}

// MockJobRouter 模拟 API 映射器
type MockJobRouter struct {
	mock.Mock
}

func (m *MockJobRouter) ExecJob(tag uint16, request iface.IRequest) error {
	args := m.Called(tag, request)
	return args.Error(0)
}

func (m *MockJobRouter) AddJob(tag uint16, job iface.IJob) iface.IJobRouter {
	m.Called(tag, job)
	return m
}

func (m *MockJobRouter) GetJob(tag uint16) iface.IJob {
	args := m.Called(tag)
	return args.Get(0).(iface.IJob)
}

// MockIRequest 模拟 IRequest 接口
type MockIRequest struct {
	mock.Mock
}

func (m *MockIRequest) Msg() iface.ISeqedTLVMsg {
	args := m.Called()
	return args.Get(0).(iface.ISeqedTLVMsg)
}

func (m *MockIRequest) Session() iface.ISession {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(iface.ISession)
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
	mockMapper := new(MockJobRouter)
	mockRequest := new(MockIRequest)
	mockMessage := new(MockIMessage)

	// Set up mock behaviors
	mockRequest.On("Msg").Return(mockMessage)
	mockRequest.On("Conn").Return(nil)
	mockMessage.On("Tag").Return("test-tag")
	mockMapper.On("ExecJob", uint16(0), mockRequest).Return(nil)

	// Initialize WorkerPool
	pool := job.NewWorkerPool(2, mockMQ, mockMapper)

	// Test Start and Stop
	t.Run("Start and Stop", func(t *testing.T) {
		mockMQ.On("Pop").Return(mockRequest).Once()
		mockMQ.On("Pop").Return(nil) // Simulate no more messages
		mockMQ.On("Close").Return().Once()

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
