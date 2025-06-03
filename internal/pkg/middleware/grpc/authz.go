package grpc

import (
	"context"

	"github.com/jwcen/miniblog/internal/pkg/contextx"
	"github.com/jwcen/miniblog/internal/pkg/errno"
	"github.com/jwcen/miniblog/internal/pkg/log"
	"google.golang.org/grpc"
)

// Authorizer 用于定义授权接口的实现.
type Authorizer interface {
	Authorize(subject, object, action string) (bool, error)
}

// AuthzInterceptor 是一个 gRPC 拦截器，用于进行请求授权.
func AuthzInterceptor(authorizer Authorizer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {

		subject := contextx.UserID(ctx) // 获取用户ID
		object := info.FullMethod       // 获取请求资源
		action := "CALL"                // 默认操作

		// 记录授权上下文信息
		log.Debugw("Build authorize context", "subject", subject, "object", object, "action", action)

		allowed, err := authorizer.Authorize(subject, object, action)
		if err != nil || !allowed {
			return nil, errno.ErrPermissionDenied.WithMessage(
				"access denied: subject=%s, object=%s, action=%s, reason=%v",
				subject,
				object,
				action,
				err,
			)
		}

		return handler(ctx, req)
	}
}
