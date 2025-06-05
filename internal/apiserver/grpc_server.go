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

package apiserver

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	handler "github.com/jwcen/miniblog/internal/apiserver/handler/grpc"
	"github.com/jwcen/miniblog/internal/pkg/log"
	mw "github.com/jwcen/miniblog/internal/pkg/middleware/grpc"
	"github.com/jwcen/miniblog/internal/pkg/server"
	apiv1 "github.com/jwcen/miniblog/pkg/api/apiserver/v1"
	genericvalidation "github.com/onexstack/onexstack/pkg/validation"
)

// grpcServer 定义一个 gRPC 服务器.
type grpcServer struct {
	srv server.Server
	// stop 为优雅关停函数.
	stop func(context.Context)
}

var _ server.Server = (*grpcServer)(nil)

// NewGRPCServerOr 创建并初始化 gRPC 或者 gRPC +  gRPC-Gateway 服务器.
// 在 Go 项目开发中，NewGRPCServerOr 这个函数命名中的 Or 一般用来表示“或者”的含义，
// 通常暗示该函数会在两种或多种选择中选择一种可能性。具体的含义需要结合函数的实现
// 或上下文来理解。以下是一些可能的解释：
//  1. 提供多种构建方式的选择
//  2. 处理默认值或回退逻辑
//  3. 表达灵活选项
func (c *ServerConfig) NewGRPCServerOr() (server.Server, error) {
	log.Infow("New GRPC Server start...")
	// 配置 gRPC 服务器选项，包括拦截器链
	serverOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			mw.RequestIDInterceptor(),
			// 认证拦截器
			selector.UnaryServerInterceptor(mw.AuthnInterceptor(c.retriever), NewAuthnWhiteListMatcher()),
			// 授权拦截器
			selector.UnaryServerInterceptor(mw.AuthzInterceptor(c.authz), NewAuthzWhiteListMatcher()),
			// Bypass 拦截器，通过所有请求的认证
			// mw.AuthnBypasswInterceptor(),
			// 为所有请求设置默认值
			mw.DefaulterInterceptor(),
			// 数据校验拦截器
			mw.ValidatorInterceptor(genericvalidation.NewValidator(c.val)),
		),
	}

	// 创建 gRPC 服务器
	grpcsrv, err := server.NewGRPCServer(
		c.cfg.GRPCOptions,
		serverOptions,
		c.cfg.TLSOptions,
		func(s grpc.ServiceRegistrar) {
			apiv1.RegisterMiniBlogServer(s, handler.NewHandler(c.biz))
		},
	)
	if err != nil {
		return nil, err
	}

	if c.cfg.ServerMode == GRPCServerMode {
		return &grpcServer{
			srv: grpcsrv,
			stop: func(ctx context.Context) {
				grpcsrv.GracefulStop(ctx)
			},
		}, nil
	}

	// 先启动 gRPC 服务器，因为 HTTP 服务器依赖 gRPC 服务器.
	go grpcsrv.RunOrDie()

	httpsrv, err := server.NewGRPCGatewayServer(
		c.cfg.HTTPOptions,
		c.cfg.GRPCOptions,
		c.cfg.TLSOptions,
		func(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
			return apiv1.RegisterMiniBlogHandler(context.Background(), mux, conn)
		},
	)
	if err != nil {
		return nil, err
	}

	return &grpcServer{
		srv: httpsrv,
		stop: func(ctx context.Context) {
			grpcsrv.GracefulStop(ctx)
			httpsrv.GracefulStop(ctx)
		},
	}, nil
}

// RunOrDie 启动 gRPC 服务器或 HTTP 反向代理服务器，异常时退出.
func (s *grpcServer) RunOrDie() {
	s.srv.RunOrDie()
}

// GracefulStop 优雅停止 HTTP 和 gRPC 服务器.
func (s *grpcServer) GracefulStop(ctx context.Context) {
	s.stop(ctx)
}

// NewAuthnWhiteListMatcher 创建认证白名单匹配器.
func NewAuthnWhiteListMatcher() selector.Matcher {
	whitelist := map[string]struct{}{
		apiv1.MiniBlog_Healthz_FullMethodName:    {},
		apiv1.MiniBlog_CreateUser_FullMethodName: {},
		apiv1.MiniBlog_Login_FullMethodName:      {},
	}

	return selector.MatchFunc(func(ctx context.Context, call interceptors.CallMeta) bool {
		_, ok := whitelist[call.FullMethod()]
		return !ok
	})
}

// NewAuthzWhiteListMatcher 创建授权白名单匹配器.
func NewAuthzWhiteListMatcher() selector.Matcher {
	whitelist := map[string]struct{}{
		apiv1.MiniBlog_Healthz_FullMethodName:    {},
		apiv1.MiniBlog_CreateUser_FullMethodName: {},
		apiv1.MiniBlog_Login_FullMethodName:      {},
	}
	return selector.MatchFunc(func(ctx context.Context, call interceptors.CallMeta) bool {
		_, ok := whitelist[call.FullMethod()]
		return !ok
	})
}
