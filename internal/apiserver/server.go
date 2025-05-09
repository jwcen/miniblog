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
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	genericoptions "github.com/onexstack/onexstack/pkg/options"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	handler "github.com/jwcen/miniblog/internal/apiserver/handler/grpc"
	"github.com/jwcen/miniblog/internal/pkg/log"
	apiv1 "github.com/jwcen/miniblog/pkg/api/apiserver/v1"
)

const (
	// GRPCServerMode 定义 gRPC 服务模式.
	// 使用 gRPC 框架启动一个 gRPC 服务器.
	GRPCServerMode = "grpc"

	// GRPCServerMode 定义 gRPC + HTTP 服务模式.
	// 使用 gRPC 框架启动一个 gRPC 服务器 + HTTP 反向代理服务器.
	GRPCGatewayServerMode = "grpc-gateway"

	// GinServerMode 定义 Gin 服务模式.
	// 使用 Gin Web 框架启动一个 HTTP 服务器.
	GinServerMode = "gin"
)

// Config 配置结构体，用于存储应用相关的配置.
// 不用 viper.Get，是因为这种方式能更加清晰的知道应用提供了哪些配置项.
type Config struct {
	ServerMode  string
	JWTKey      string
	Expiration  time.Duration
	GRPCOptions *genericoptions.GRPCOptions
}

// UnionServer 定义一个联合服务器. 根据 ServerMode 决定要启动的服务器类型.
type UnionServer struct {
	cfg *Config
	srv *grpc.Server
	lis net.Listener
}

// NewUnionServer 根据配置创建联合服务器.
func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	lis, err := net.Listen("tcp", cfg.GRPCOptions.Addr)
	if err != nil {
		log.Fatalw("Failed to listen", "err", err)
		return nil, err
	}

	grpcsrv := grpc.NewServer()
	apiv1.RegisterMiniBlogServer(grpcsrv, handler.NewHandler())
	reflection.Register(grpcsrv)

	return &UnionServer{
		cfg: cfg,
		srv: grpcsrv,
		lis: lis,
	}, nil
}

// Run 运行应用.
func (s *UnionServer) Run() error {
	log.Infow("ServerMode from ServerOptions", "jwt-key", s.cfg.JWTKey)
	log.Infow("ServerMode from Viper", "jwt-key", viper.GetString("jwt-key"))

	jsonData, _ := json.MarshalIndent(s.cfg, "", "  ")
	fmt.Println(string(jsonData))

	// 启动 gRPC 服务器
	go func() {
		log.Infow("Starting gRPC server", "addr", s.cfg.GRPCOptions.Addr)
		if err := s.srv.Serve(s.lis); err != nil {
			log.Fatalw("Failed to serve", "err", err)
		}
	}()

	// 等待中断信号优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	// 监听 SIGINT 和 SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
	s.srv.GracefulStop() // 优雅关闭 gRPC 服务器
	return nil
}
