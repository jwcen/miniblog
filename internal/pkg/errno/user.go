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

package errno

import (
	"net/http"

	"github.com/onexstack/onexstack/pkg/errorsx"
)

var (
	// ErrUsernameInvalid 表示用户名不合法.
	ErrUsernameInvalid = &errorsx.ErrorX{
		Code:    http.StatusBadRequest,
		Reason:  "InvalidArgument.UsernameInvalid",
		Message: "Invalid username: Username must consist of letters, digits, and underscores only, and its length must be between 3 and 20 characters.",
	}

	// ErrPasswordInvalid 表示密码不合法.
	ErrPasswordInvalid = &errorsx.ErrorX{
		Code:    http.StatusBadRequest,
		Reason:  "InvalidArgument.PasswordInvalid",
		Message: "Password is incorrect.",
	}

	// ErrUserAlreadyExists 表示用户已存在.
	ErrUserAlreadyExists = &errorsx.ErrorX{Code: http.StatusBadRequest, Reason: "AlreadyExist.UserAlreadyExists", Message: "User already exists."}

	// ErrUserNotFound 表示未找到指定用户.
	ErrUserNotFound = &errorsx.ErrorX{Code: http.StatusNotFound, Reason: "NotFound.UserNotFound", Message: "User not found."}
)
