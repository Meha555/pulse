package job

import (
	"bytes"
	"log"
	"regexp"
	"testing"

	"github.com/Meha555/pulse/core/message"
	"github.com/Meha555/pulse/server/common"

	"github.com/Meha555/go-tinylog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMsgQueue 模拟消息队列
type MockMsgQueue struct {
	mock.Mock
}

func (m *MockMsgQueue) Push(request common.IRequest) {
	m.Called(request)
}

func (m *MockMsgQueue) Pop() common.IRequest {
	args := m.Called()
	return args.Get(0).(common.IRequest)
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

func (m *MockJobRouter) ExecJob(tag uint16, request common.IRequest) error {
	args := m.Called(tag, request)
	return args.Error(0)
}

func (m *MockJobRouter) AddJob(tag uint16, job IJob) IJobRouter {
	m.Called(tag, job)
	return m
}

func (m *MockJobRouter) GetJob(tag uint16) IJob {
	args := m.Called(tag)
	return args.Get(0).(IJob)
}

// MockIRequest 模拟 IRequest 接口
type MockIRequest struct {
	mock.Mock
}

func (m *MockIRequest) Msg() message.ISeqedTLVMsg {
	args := m.Called()
	return args.Get(0).(message.ISeqedTLVMsg)
}

func (m *MockIRequest) Session() common.ISession {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(common.ISession)
}

func (m *MockIRequest) Get(key string) (value interface{}, exists bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *MockIRequest) Set(key string, value interface{}) {
	m.Called(key, value)
}

// MockISeqedTLVMsg 模拟 ISeqedTLVMsg 接口
type MockISeqedTLVMsg struct {
	mock.Mock
}

func (m *MockISeqedTLVMsg) Tag() uint16 {
	args := m.Called()
	return args.Get(0).(uint16)
}

func (m *MockISeqedTLVMsg) Serial() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

func (m *MockISeqedTLVMsg) SetSerial(seq uint32) {
	m.Called(seq)
}

func (m *MockISeqedTLVMsg) Body() []byte {
	return m.Called().Get(0).([]byte)
}

func (m *MockISeqedTLVMsg) BodyLen() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

func (m *MockISeqedTLVMsg) HeaderLen() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

func (m *MockISeqedTLVMsg) SetBody(body []byte) {
	m.Called(body)
}

func (m *MockISeqedTLVMsg) SetTag(tag uint16) {
	m.Called(tag)
}

func TestWorkerPool(t *testing.T) {
	// Redirect log output to buffer for verification
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(nil)
	logger.SetLevel(tinylog.LevelDebug)

	// Create mocks
	mockMQ := new(MockMsgQueue)
	mockMapper := new(MockJobRouter)
	mockRequest := new(MockIRequest)
	mockMessage := new(MockISeqedTLVMsg)

	var testTag uint16 = 1

	// Set up mock behaviors
	mockRequest.On("Msg").Return(mockMessage)
	mockRequest.On("Conn").Return(nil)
	mockMessage.On("Tag").Return(testTag)
	mockMapper.On("ExecJob", uint16(0), mockRequest).Return(nil)

	// Initialize WorkerPool
	pool := NewWorkerPool(2, mockMQ, mockMapper)

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

	// Test Post with nil request
	t.Run("Post with nil request", func(t *testing.T) {
		mockMQ.On("Push", nil).Return()
		pool.Post(nil)
		mockMQ.AssertNumberOfCalls(t, "Push", 0)
	})

	// Test Start with zero workers
	t.Run("Start with zero workers", func(t *testing.T) {
		poolZero := NewWorkerPool(0, mockMQ, mockMapper)
		poolZero.Start()
		poolZero.Stop()
		assert.Contains(t, logBuffer.String(), "All workers stopped")
	})

	// Test ExecJob error
	t.Run("ExecJob error", func(t *testing.T) {
		mockMapper.On("ExecJob", uint16(0), mockRequest).Return(assert.AnError)
		mockMQ.On("Pop").Return(mockRequest).Once()
		mockMQ.On("Pop").Return(nil)
		mockMQ.On("Close").Return().Once()

		pool.Start()
		pool.Stop()

		assert.Regexp(t, regexp.MustCompile(`Worker\[\d+\] process request failed`), logBuffer.String())
	})
}
