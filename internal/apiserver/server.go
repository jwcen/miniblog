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
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jwcen/miniblog/internal/pkg/known"
	mw "github.com/jwcen/miniblog/internal/pkg/middleware/gin"
	"github.com/onexstack/onexstack/pkg/authz"
	genericoptions "github.com/onexstack/onexstack/pkg/options"
	"github.com/onexstack/onexstack/pkg/store/where"
	"github.com/onexstack/onexstack/pkg/token"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"k8s.io/utils/ptr"

	"github.com/jwcen/miniblog/internal/apiserver/biz"
	"github.com/jwcen/miniblog/internal/apiserver/model"
	"github.com/jwcen/miniblog/internal/apiserver/pkg/validation"
	"github.com/jwcen/miniblog/internal/apiserver/store"
	"github.com/jwcen/miniblog/internal/pkg/contextx"
	"github.com/jwcen/miniblog/internal/pkg/log"
	"github.com/jwcen/miniblog/internal/pkg/server"
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
	ServerMode   string
	JWTKey       string
	Expiration   time.Duration
	GRPCOptions  *genericoptions.GRPCOptions
	HTTPOptions  *genericoptions.HTTPOptions
	TLSOptions   *genericoptions.TLSOptions
	MySQLOptions *genericoptions.MySQLOptions
	EnableMemoryStore bool
}

// UnionServer 定义一个联合服务器. 根据 ServerMode 决定要启动的服务器类型.
//
// 联合服务器分为以下 2 大类：
//  1. Gin 服务器：由 Gin 框架创建的标准的 REST 服务器。根据是否开启 TLS，
//     来判断启动 HTTP 或者 HTTPS；
//  2. GRPC 服务器：由 gRPC 框架创建的标准 RPC 服务器
//  3. HTTP 反向代理服务器：由 grpc-gateway 框架创建的 HTTP 反向代理服务器。
//     根据是否开启 TLS，来判断启动 HTTP 或者 HTTPS；
//
// HTTP 反向代理服务器依赖 gRPC 服务器，所以在开启 HTTP 反向代理服务器时，会先启动 gRPC 服务器.
type UnionServer struct {
	srv server.Server
}

// ServerConfig 包含服务器的核心依赖和配置.
type ServerConfig struct {
	cfg       *Config
	biz       biz.IBiz
	val       *validation.Validator
	retriever mw.UserRetriever
	authz     *authz.Authz
}

// NewServerConfig 创建一个 *ServerConfig 实例.
// 进阶：这里其实可以使用依赖注入的方式，来创建 *ServerConfig.
func (cfg *Config) NewServerConfig() (*ServerConfig, error) {
	db, err := cfg.NewDB()
	if err != nil {
		log.Errorw("cfg.NewDB() err:", err)
		return nil, err
	}

	store := store.NewStore(db)

	authz, err := authz.NewAuthz(store.DB(context.TODO()))
	if err != nil {
		return nil, err
	}

	return &ServerConfig{
		cfg:       cfg,
		biz:       biz.NewBiz(store, authz),
		val:       validation.New(store),
		retriever: &UserRetriever{store},
		authz:     authz,
	}, nil
}

func (cfg *Config) NewDB() (*gorm.DB, error) {
	if !cfg.EnableMemoryStore {
		log.Infow("Initializing database connection", "type", "mysql", "addr", cfg.MySQLOptions.Addr)
        return cfg.MySQLOptions.NewDB()
	}

	log.Infow("Initializing database connection", "type", "memory", "engine", "SQLite")

	// 使用SQLite内存模式配置数据库
    // ?cache=shared 用于设置 SQLite 的缓存模式为 共享缓存模式 (shared)。
    // 默认情况下，SQLite 的每个数据库连接拥有自己的独立缓存，这种模式称为 专用缓存 (private)。
    // 使用 共享缓存模式 (shared) 后，不同连接可以共享同一个内存中的数据库和缓存。
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        log.Errorw("Failed to create database connection", "err", err)
        return nil, err
    }

	// 自动迁移数据库结构
    if err := db.AutoMigrate(&model.UserM{}, &model.PostM{}, &model.CasbinRuleM{}); err != nil {
        log.Errorw("Failed to migrate database schema", "err", err)
        return nil, err
    }

	// 插入 casbin_rule 表记录
    adminR, userR := "role::admin", "role::user"
    casbinRules := []model.CasbinRuleM{
        {PType: ptr.To("g"), V0: ptr.To("user-000000"), V1: &adminR},
        {PType: ptr.To("p"), V0: &adminR, V1: ptr.To("*"), V2: ptr.To("*"), V3: ptr.To("allow")},
        {PType: ptr.To("p"), V0: &userR, V1: ptr.To("/v1.MiniBlog/DeleteUser"), V2: ptr.To("CALL"), V3: ptr.To("deny")},
        {PType: ptr.To("p"), V0: &userR, V1: ptr.To("/v1.MiniBlog/ListUser"), V2: ptr.To("CALL"), V3: ptr.To("deny")},
        {PType: ptr.To("p"), V0: &userR, V1: ptr.To("/v1/users"), V2: ptr.To("GET"), V3: ptr.To("deny")},
        {PType: ptr.To("p"), V0: &userR, V1: ptr.To("/v1/users/*"), V2: ptr.To("DELETE"), V3: ptr.To("deny")},
    }

    if err := db.Create(&casbinRules).Error; err != nil {
        log.Fatalw("Failed to insert casbin_rule records", "err", err)
        return nil, err
    }

    // 插入默认用户（root用户）
    user := model.UserM{
        UserID:    "user-000000",
        Username:  "root",
        Password:  "miniblog1234",
        Nickname:  "administrator",
        Email:     "colin404@foxmail.com",
        Phone:     "18110000000",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := db.Create(&user).Error; err != nil {
        log.Fatalw("Failed to insert default root user", "err", err)
        return nil, err
    }

    return db, nil
}

// NewUnionServer 根据配置创建联合服务器.
func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	// 注册租户解析函数，通过上下文获取用户 ID
	//nolint: gocritic
	where.RegisterTenant("userID", func(ctx context.Context) string {
		return contextx.UserID(ctx)
	})

	token.Init(cfg.JWTKey, known.XUserID, cfg.Expiration)
	
	log.Infow("Initializing federation server", "server-mode", cfg.ServerMode, "enable-memory-store", cfg.EnableMemoryStore)

	srv, err := InitializeWebServer(cfg)
	if err != nil {
		return nil, err
	}

	log.Infow("Initializing federation server", "server-mode", cfg.ServerMode)

	return &UnionServer{srv: srv}, nil
}

// Run 运行应用.
func (s *UnionServer) Run() error {
	s.srv.RunOrDie()

	quit := make(chan os.Signal, 1)
	// / 当执行 kill 命令时（不带参数），默认会发送 syscall.SIGTERM 信号
	// 使用 kill -2 命令会发送 syscall.SIGINT 信号（例如按 CTRL+C 触发）
	// 使用 kill -9 命令会发送 syscall.SIGKILL 信号，但 SIGKILL 信号无法被捕获，因此无需监听和处理
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Infow("Shutting down server ...")
	// 优雅关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 先关闭依赖的服务，再关闭被依赖的服务
	s.srv.GracefulStop(ctx)

	log.Infow("Server exited")

	return nil
}

// UserRetriever 定义一个用户数据获取器. 用来获取用户信息.
type UserRetriever struct {
	store store.IStore
}

func (u *UserRetriever) GetUser(ctx context.Context, userID string) (*model.UserM, error) {
	return u.store.User().Get(ctx, where.F("userID", userID))
}

// ProvideDB 根据配置提供一个数据库实例。
func ProvideDB(cfg *Config) (*gorm.DB, error) {
	return cfg.NewDB()
}

func NewWebServer(serverMode string, serverConfig *ServerConfig) (server.Server, error) {
	// 根据服务模式创建对应的服务实例
	// 实际企业开发中，可以根据需要只选择一种服务器模式.
	// 这里为了方便展示，通过 cfg.ServerMode 同时支持了 Gin 和 GRPC 2 种服务器模式.
	// 默认为 gRPC 服务器模式.
	switch serverMode {
	case GinServerMode:
		return serverConfig.NewGinServer(), nil
	default:
		return serverConfig.NewGRPCServerOr()
	}
}
