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
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	handler "github.com/jwcen/miniblog/internal/apiserver/handler/http"
	mw "github.com/jwcen/miniblog/internal/pkg/middleware/gin"
	"github.com/jwcen/miniblog/internal/pkg/server"
)

type ginServer struct {
	srv server.Server
}

var _ server.Server = (*ginServer)(nil)

func (c *ServerConfig) NewGinServer() server.Server {
	engine := gin.New()
	// 注册全局中间件，用于恢复 panic、设置 HTTP 头、添加请求 ID 等
	engine.Use(gin.Recovery(), mw.NoCache, mw.Cors, mw.Secure, mw.RequestIDMiddleware())

	// 注册 REST API 路由
	c.InstallRESTAPI(engine)

	return &ginServer{srv: nil}
}

// RunOrDie 启动 Gin 服务器，出错则程序崩溃退出.
func (s *ginServer) RunOrDie() {
	select {}
}

// GracefulStop 优雅停止服务器.
func (s *ginServer) GracefulStop(ctx context.Context) {}

// 注册 API 路由。路由的路径和 HTTP 方法，严格遵循 REST 规范.
func (c *ServerConfig) InstallRESTAPI(engine *gin.Engine) {
	// 注册业务无关的 API 接口
	InstallGenericAPI(engine)

	// 创建核心业务处理器
	handler := handler.NewHandler()

	// 注册健康检查接口
	engine.GET("/healthz", handler.Healthz)
}

// InstallGenericAPI 注册业务无关的路由，例如 pprof、404 处理等.
func InstallGenericAPI(engine *gin.Engine) {
	// 注册 pprof 路由
	pprof.Register(engine)

	// 注册 404 路由处理
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "Page not found.")
	})
}
