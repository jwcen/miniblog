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


package rid_test

import (
	"fmt"

	"github.com/jwcen/miniblog/internal/pkg/rid"
)

func ExampleResourceID_String() {
	// 定义一个资源标识符，例如用户资源
	userID := rid.UserID

	// 调用String方法，将ResourceID类型转换为字符串类型
	idString := userID.String()

	// 输出结果
	fmt.Println(idString)

	// Output:
	// user
}
