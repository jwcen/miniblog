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

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jwcen/miniblog/internal/pkg/errno"
	apiv1 "github.com/jwcen/miniblog/pkg/api/apiserver/v1"
	genericvalidation "github.com/onexstack/onexstack/pkg/validation"
)

// Validate 校验字段的有效性.
func (v *Validator) ValidatePostRules() genericvalidation.Rules {
	// 定义各字段的校验逻辑，通过一个 map 实现模块化和简化
	return genericvalidation.Rules{
		"PostID": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("postID cannot be empty")
			}
			return nil
		},
		"Title": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("title cannot be empty")
			}
			return nil
		},
		"Content": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("content cannot be empty")
			}
			return nil
		},
	}
}

// ValidateCreatePostRequest 校验 CreatePostRequest 结构体的有效性.
func (v *Validator) ValidateCreatePostRequest(ctx context.Context, rq *apiv1.CreatePostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

// ValidateUpdatePostRequest 校验更新用户请求.
func (v *Validator) ValidateUpdatePostRequest(ctx context.Context, rq *apiv1.UpdatePostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

// ValidateDeletePostRequest 校验 DeletePostRequest 结构体的有效性.
func (v *Validator) ValidateDeletePostRequest(ctx context.Context, rq *apiv1.DeletePostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

// ValidateGetPostRequest 校验 GetPostRequest 结构体的有效性.
func (v *Validator) ValidateGetPostRequest(ctx context.Context, rq *apiv1.GetPostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

// ValidateListPostRequest 校验 ListPostRequest 结构体的有效性.
func (v *Validator) ValidateListPostRequest(ctx context.Context, rq *apiv1.ListPostRequest) error {
	if err := validation.Validate(rq.GetTitle(), validation.Length(5, 100), is.URL); err != nil {
		return errno.ErrInvalidArgument.WithMessage(err.Error())
	}
	return genericvalidation.ValidateSelectedFields(rq, v.ValidatePostRules(), "Offset", "Limit")
}
