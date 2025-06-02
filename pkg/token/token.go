package token

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	key         string
	identityKey string
	expiration  time.Duration
}

var (
	config = Config{"Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", "identityKey", 2 * time.Hour}
	once   sync.Once
)

// Init 设置包级别的配置 config, config 会用于本包后面的 token 签发和解析.
func Init(key, identityKey string, expiration time.Duration) {
	once.Do(func() {
		if key != "" {
			config.key = key
		}

		if identityKey != "" {
			config.identityKey = identityKey
		}

		if expiration != 0 {
			config.expiration = expiration
		}
	})
}

// Parse 使用指定的密钥 key 解析 token，解析成功返回 token 上下文，否则报错.
func Parse(tokenString, key string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保 token 加密算法是预期的加密算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(key), nil
	})
	if err != nil {
		return "", err
	}

	// 如果解析成功，从 token 中取出 token 的主题
	var identityKey string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if key, exists := claims[config.identityKey]; exists {
			if identity, valid := key.(string); valid {
				identityKey = identity
			}
		}
	}

	if identityKey == "" {
		return "", jwt.ErrSignatureInvalid
	}

	return identityKey, nil
}

// ParseRequest 从请求头中获取令牌，并将其传递给 Parse 函数以解析令牌.
func ParseRequest(ctx context.Context) (string, error) {
	var (
		token string
		err   error
	)

	switch typed := ctx.(type) {
	case *gin.Context:
		header := typed.Request.Header.Get("Authorization")
		if len(header) == 0 {
			return "", errors.New("the length of the `Authorization` header is zero")
		}

		_, _ = fmt.Sscanf(header, "Bearer %s", &token)

	default:
		// gRPC 服务
		token, err = auth.AuthFromMD(typed, "Bearer")
		if err != nil {
			return "", status.Errorf(codes.Unauthenticated, "invalid auth token")
		}
	}

	return Parse(token, config.key) // 解析 token
}

// Sign 使用 jwtSecret 签发 token，token 的 claims 中会存放传入的 subject.
func Sign(identityKey string) (string, time.Time, error) {
	expireAt := time.Now().Add(config.expiration)

	// Token 的内容
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.identityKey: identityKey,       // 存放用户身份
		"nbf":              time.Now().Unix(), // token 生效时间
		"iat":              time.Now().Unix(), // token 签发时间
		"exp":              expireAt.Unix(),   // token 过期时间
	})

	// 签发 token
	tokenString, err := token.SignedString([]byte(config.key))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expireAt, nil
}
