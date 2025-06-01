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

package validation

import (
	"regexp"

	"github.com/jwcen/miniblog/internal/apiserver/store"
	"github.com/jwcen/miniblog/internal/pkg/errno"
)

type Validator struct {
	store store.IStore
}

func New(store store.IStore) *Validator {
	return &Validator{store: store}
}

// isValidUsername 校验用户名是否合法.
func isValidUsername(username string) bool {
	// 用户名必须仅包含字母、数字和下划线，并且长度在 3 到 20 个字符之间
	var (
		lengthRegex = `^.{3,20}$`       // 长度在 3 到 20 个字符之间
		validRegex  = `^[A-Za-z0-9_]+$` // 仅包含字母、数字和下划线
	)

	// 校验长度
	if matched, _ := regexp.MatchString(lengthRegex, username); !matched {
		return false
	}
	// 校验字符合法性
	if matched, _ := regexp.MatchString(validRegex, username); !matched {
		return false
	}
	return true
}

// isValidPassword 判断密码是否符合复杂度要求.
func isValidPassword(password string) error {
	// 检查新密码是否为空
	if password == "" {
		return errno.ErrInvalidArgument.WithMessage("password cannot be empty")
	}

	// 检查新密码的长度要求
	if len(password) < 6 {
		return errno.ErrInvalidArgument.WithMessage("password must be at least 6 characters long")
	}

	// 使用正则表达式检查是否至少包含一个字母
	letterPattern := regexp.MustCompile(`[A-Za-z]`)
	if !letterPattern.MatchString(password) {
		return errno.ErrInvalidArgument.WithMessage("password must contain at least one letter")
	}

	// 使用正则表达式检查是否至少包含一个数字
	numberPattern := regexp.MustCompile(`\d`)
	if !numberPattern.MatchString(password) {
		return errno.ErrInvalidArgument.WithMessage("password must contain at least one number")
	}

	return nil
}

// isValidEmail 判断电子邮件是否合法.
func isValidEmail(email string) error {
	// 检查电子邮件地址格式
	if email == "" {
		return errno.ErrInvalidArgument.WithMessage("email cannot be empty")
	}

	// 使用正则表达式校验电子邮件格式
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailPattern.MatchString(email) {
		return errno.ErrInvalidArgument.WithMessage("invalid email format")
	}

	return nil
}

// isValidPhone 判断手机号码是否合法.
func isValidPhone(phone string) error {
	// 检查手机号码格式
	if phone == "" {
		return errno.ErrInvalidArgument.WithMessage("phone cannot be empty")
	}

	// 使用正则表达式校验手机号码格式（假设是中国手机号，11位数字）
	phonePattern := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phonePattern.MatchString(phone) {
		return errno.ErrInvalidArgument.WithMessage("invalid phone format")
	}

	return nil
}
