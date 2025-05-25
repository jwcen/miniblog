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

package biz

import (
	postV1 "github.com/jwcen/miniblog/internal/apiserver/biz/v1/post"
	userV1 "github.com/jwcen/miniblog/internal/apiserver/biz/v1/user"
)

type IBiz interface {
	// 获取用户业务接口.
	UserV1() userV1.UserBiz
	// 获取帖子业务接口.
	POSTV1() postV1.PostBiz
}
