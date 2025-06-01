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
	"context"

	genericvalidation "github.com/onexstack/onexstack/pkg/validation"

	"github.com/jwcen/miniblog/internal/pkg/contextx"
	"github.com/jwcen/miniblog/internal/pkg/errno"
	apiv1 "github.com/jwcen/miniblog/pkg/api/apiserver/v1"
)

func (v *Validator) ValidateUserRules() genericvalidation.Rules {
	// 通用的密码校验函数
	validatePassword := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			return isValidPassword(value.(string))
		}
	}

	// 定义各字段的校验逻辑，通过一个 map 实现模块化和简化
	return genericvalidation.Rules{
		"Password":    validatePassword(),
		"OldPassword": validatePassword(),
		"NewPassword": validatePassword(),
		"UserID": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("userID cannot be empty")
			}
			return nil
		},
		"Username": func(value any) error {
			if !isValidUsername(value.(string)) {
				return errno.ErrUsernameInvalid
			}
			return nil
		},
		"Nickname": func(value any) error {
			if len(value.(string)) >= 30 {
				return errno.ErrInvalidArgument.WithMessage("nickname must be less than 30 characters")
			}
			return nil
		},
		"Email": func(value any) error {
			return isValidEmail(value.(string))
		},
		"Phone": func(value any) error {
			return isValidPhone(value.(string))
		},
		"Limit": func(value any) error {
			if value.(int64) <= 0 {
				return errno.ErrInvalidArgument.WithMessage("limit must be greater than 0")
			}
			return nil
		},
		"Offset": func(value any) error {
			return nil
		},
	}
}

// ValidateLogin 校验修改密码请求.
func (v *Validator) ValidateLoginRequest(ctx context.Context, rq *apiv1.LoginRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

// ValidateChangePasswordRequest 校验 ChangePasswordRequest 结构体的有效性.
func (v *Validator) ValidateChangePasswordRequest(ctx context.Context, rq *apiv1.ChangePasswordRequest) error {
	if rq.GetUserID() != contextx.UserID(ctx) {
		return errno.ErrPermissionDenied.WithMessage("The logged-in user `%s` does not match request user `%s`", contextx.UserID(ctx), rq.GetUserID())
	}
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

// ValidateCreateUserRequest 校验 CreateUserRequest 结构体的有效性.
func (v *Validator) ValidateCreateUserRequest(ctx context.Context, rq *apiv1.CreateUserRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

// ValidateUpdateUserRequest 校验更新用户请求.
func (v *Validator) ValidateUpdateUserRequest(ctx context.Context, rq *apiv1.UpdateUserRequest) error {
	if rq.GetUserID() != contextx.UserID(ctx) {
		return errno.ErrPermissionDenied.WithMessage("The logged-in user `%s` does not match request user `%s`", contextx.UserID(ctx), rq.GetUserID())
	}
	return genericvalidation.ValidateSelectedFields(rq, v.ValidateUserRules(), "UserID")
}

// ValidateDeleteUserRequest 校验 DeleteUserRequest 结构体的有效性.
func (v *Validator) ValidateDeleteUserRequest(ctx context.Context, rq *apiv1.DeleteUserRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

// ValidateGetUserRequest 校验 GetUserRequest 结构体的有效性.
func (v *Validator) ValidateGetUserRequest(ctx context.Context, rq *apiv1.GetUserRequest) error {
	if rq.GetUserID() != contextx.UserID(ctx) {
		return errno.ErrPermissionDenied.WithMessage("The logged-in user `%s` does not match request user `%s`", contextx.UserID(ctx), rq.GetUserID())
	}
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

// ValidateListUserRequest 校验 ListUserRequest 结构体的有效性.
func (v *Validator) ValidateListUserRequest(ctx context.Context, rq *apiv1.ListUserRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}
