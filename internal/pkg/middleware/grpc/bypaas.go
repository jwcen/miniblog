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

	"github.com/jwcen/miniblog/internal/pkg/contextx"
	"github.com/jwcen/miniblog/internal/pkg/known"
	"github.com/jwcen/miniblog/internal/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthnBypasswInterceptor 是一个 gRPC 拦截器，模拟所有请求都通过认证。
func AuthnBypasswInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {

		userID := "user-000001"
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get(known.XUserID); len(values) > 0 {
				userID = values[0]
			}
		}

		log.Debugw("Simulated authentication successful", "userID", userID)

		// 将默认的用户信息存入上下文
		//nolint: staticcheck
		ctx = context.WithValue(ctx, known.XUserID, userID)
		// 为 log 和 contextx 提供用户上下文支持
		ctx = contextx.WithUserID(ctx, userID)

		return handler(ctx, req)
	}
}
