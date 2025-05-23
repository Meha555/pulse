package utils_test

import (
	"testing"

	"my-zinx/utils"
)

func TestDict(t *testing.T) {
	t.Run("Store", func(t *testing.T) {
		dict := utils.Dict[string, int]{}

		// 测试存储单个值
		dict.Store("key1", 10)
		value, exists := dict.Load("key1")
		if !exists || value != 10 {
			t.Errorf("Expected value 10 for key1, got %v", value)
		}

		// 测试覆盖已存在的键
		dict.Store("key1", 20)
		value, exists = dict.Load("key1")
		if !exists || value != 20 {
			t.Errorf("Expected value 20 for key1 after overwrite, got %v", value)
		}
	})

	t.Run("Load", func(t *testing.T) {
		dict := utils.Dict[string, int]{}
		dict.Store("key1", 10)

		// 测试加载存在的键
		value, exists := dict.Load("key1")
		if !exists || value != 10 {
			t.Errorf("Expected value 10 for key1, got %v", value)
		}

		// 测试加载不存在的键
		_, exists = dict.Load("nonexistent")
		if exists {
			t.Error("Expected nonexistent key to return false")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		dict := utils.Dict[string, int]{}
		dict.Store("key1", 10)

		// 测试删除存在的键
		dict.Delete("key1")
		_, exists := dict.Load("key1")
		if exists {
			t.Error("Expected key1 to be deleted")
		}

		// 测试删除不存在的键
		dict.Delete("nonexistent")
		_, exists = dict.Load("nonexistent")
		if exists {
			t.Error("Expected nonexistent key to remain deleted")
		}
	})
}
