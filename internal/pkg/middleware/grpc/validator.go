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

package grpc

import (
	"context"

	"google.golang.org/grpc"
)

// RequestValidator 定义了用于自定义验证的接口.
type RequestValidator interface {
	Validate(ctx context.Context, rq any) error
}

// ValidatorInterceptor 是一个 gRPC 拦截器，用于对请求进行验证.
func ValidatorInterceptor(validator RequestValidator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, rq any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// 调用自定义验证方法
		if err := validator.Validate(ctx, rq); err != nil {
			// 注意这里不用返回 errno.ErrInvalidArgument 类型的错误信息，由 validator.Validate 返回.
			return nil, err // 返回验证错误
		}

		// 继续处理请求
		return handler(ctx, rq)
	}
}
