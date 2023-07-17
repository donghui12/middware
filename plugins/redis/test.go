package redis

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestRedisSetAndGet(t *testing.T) {
	conf := &Config{
		Host:     "localhost", // 根据实际情况修改Redis的主机地址
		Port:     6379,        // 根据实际情况修改Redis的端口号
		Username: "",          // 根据实际情况修改Redis的用户名（如果有）
		Password: "",          // 根据实际情况修改Redis的密码（如果有）
		DB:       0,           // 根据实际情况修改Redis的数据库索引
	}

	err := Init(conf)
	if err != nil {
		t.Errorf("Failed to initialize Redis client: %v", err)
		return
	}

	ctx := context.TODO()

	// 测试Set函数
	key := "test_key"
	value := map[string]interface{}{
		"name":  "John Doe",
		"email": "john.doe@example.com",
	}
	expiration := 5 * time.Minute

	err = Client().Set(ctx, key, value, expiration)
	if err != nil {
		t.Errorf("Failed to set value to Redis: %v", err)
		return
	}

	// 测试Get函数
	var retrievedValue map[string]interface{}
	err = Client().Get(ctx, key, &retrievedValue)
	if err != nil {
		t.Errorf("Failed to get value from Redis: %v", err)
		return
	}

	// 检查获取到的值是否与设置的值相等
	if !isEqual(value, retrievedValue) {
		t.Errorf("Retrieved value does not match the expected value")
		return
	}
}

// 判断两个map是否相等的辅助函数
func isEqual(a, b map[string]interface{}) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}
