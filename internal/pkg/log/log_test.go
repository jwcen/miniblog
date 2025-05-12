// Copyright 2024 jayvee <jvvcen@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockLogger 用于测试的自定义 Logger
var mockLogger *zapLogger

// TestMain 初始化测试环境
func TestMain(m *testing.M) {
	opts := &Options{
		Level:             "debug",
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            "json",
		OutputPaths:       []string{"stdout"},
	}
	Init(opts)
	mockLogger = std
	m.Run()
}

// TestLoggerMethods 测试日志记录方法
func TestLoggerMethods(t *testing.T) {
	// 测试不同级别的日志记录
	assert.NotPanics(t, func() {
		Debugw("debug message", "key1", "value1")
		Infow("info message", "key2", "value2")
		Warnw("warn message", "key3", "value3")
		Errorw("error message", "key4", "value4")
	}, "Log methods should not cause a panic in this test")

	assert.Panics(t, func() {
		Panicw("fatal message", "key6", "value6") // 这会导致程序退出
	}, "Panicw should cause a panic and exit the program")
}

// TestLoggerInitialization 测试 Logger 初始化
func TestLoggerInitialization(t *testing.T) {
	opts := NewOptions()
	logger := New(opts)

	assert.NotNil(t, logger, "Logger should not be nil after initialization")
	assert.IsType(t, &zapLogger{}, logger, "Logger should be of type *zapLogger")
}

// TestSync 测试日志同步
func TestSync(t *testing.T) {
	assert.NotPanics(t, func() {
		Sync() // 确保 Sync 不会引发恐慌
	}, "Sync should not panic")
}
