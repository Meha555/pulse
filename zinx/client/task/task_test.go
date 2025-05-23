package task

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskHandler struct {
	mock.Mock
}

func (m *MockTaskHandler) Do(t *Task) error {
	args := m.Called(t)
	return args.Error(0)
}

func TestTask_NewTask(t *testing.T) {
	id := uuid.New()
	t.Logf("Task %s", id)
	mockHandler := new(MockTaskHandler)
	task := NewTask(id, mockHandler.Do, WithData("hello"))
	task.AppendData("world")

	assert.Equal(t, task.Data(), []interface{}{"hello", "world"})
}

func TestTask_Exec(t *testing.T) {
	id := uuid.New()
	t.Logf("Task %s", id)
	mockHandler := new(MockTaskHandler)
	task := NewTask(id, mockHandler.Do)

	var expectedErr = errors.New("mock error")
	mockHandler.On("Do", task).Return(expectedErr).Once()
	task.Exec()
	<-task.Done()
	assert.Equal(t, TaskStatusFinished, task.Status())
	if assert.Error(t, task.Err()) {
		assert.Equal(t, expectedErr, task.Err())
	}
}

func TestTask_Exec_Timeout(t *testing.T) {
	id := uuid.New()
	t.Logf("Task %s", id)
	mockHandler := new(MockTaskHandler)
	task := NewTask(id, mockHandler.Do, WithTimeout(1*time.Millisecond))

	mockHandler.On("Do", task).Return(nil).Once()
	time.Sleep(2 * time.Millisecond)
	task.Exec()
	<-task.Done()
	assert.Equal(t, TaskStatusCanceled, task.Status())
	assert.NoError(t, task.Err())
}

func TestTask_Exec_Panic(t *testing.T) {
	id := uuid.New()
	t.Logf("Task %s", id)
	mockHandler := new(MockTaskHandler)
	task := NewTask(id, mockHandler.Do)

	mockHandler.On("Do", task).Return(nil).Panic("mock panic").Once()
	task.Exec()
	<-task.Done()
	assert.Equal(t, TaskStatusFailed, task.Status())
	if assert.Error(t, task.Err()) {
		assert.Regexp(t, fmt.Sprintf(`^task %s panic: mock panic$`, id), task.Err().Error())
	}
}

func TestTask_Exec_Cancel(t *testing.T) {
	id := uuid.New()
	t.Logf("Task %s", id)
	mockHandler := new(MockTaskHandler)
	task := NewTask(id, mockHandler.Do)

	mockHandler.On("Do", task).Return(nil).After(2 * time.Millisecond).Once()
	task.Exec()
	task.Cancel()
	<-task.Done()
	assert.Equal(t, TaskStatusCanceled, task.Status())
	assert.NoError(t, task.Err())
}

func TestTask_Exec_Repeat(t *testing.T) {
	id := uuid.New()
	t.Logf("Task %s", id)
	mockHandler := new(MockTaskHandler)
	task := NewTask(id, mockHandler.Do, WithRepeat(2))

	mockHandler.On("Do", task).Return(nil).Twice()
	task.Exec()
	<-task.Done()
	mockHandler.AssertNumberOfCalls(t, "Do", 2)
	assert.Equal(t, TaskStatusFinished, task.Status())
	assert.NoError(t, task.Err())
}
