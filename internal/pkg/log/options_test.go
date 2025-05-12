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

func TestNewOptions(t *testing.T) {
	// 调用 NewOptions 函数创建新的 Options 实例
	opts := NewOptions()

	// 验证 Options 的默认值
	assert.NotNil(t, opts, "Options should not be nil")
	assert.Equal(t, false, opts.DisableCaller, "DisableCaller should be false by default")
	assert.Equal(t, false, opts.DisableStacktrace, "DisableStacktrace should be false by default")
	assert.Equal(t, "info", opts.Level, "Level should be 'info' by default")
	assert.Equal(t, "console", opts.Format, "Format should be 'console' by default")
	assert.Equal(t, []string{"stdout"}, opts.OutputPaths, "OutputPaths should be ['stdout'] by default")
}
