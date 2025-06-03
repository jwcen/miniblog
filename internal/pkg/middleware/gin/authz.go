package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jwcen/miniblog/internal/pkg/contextx"
	"github.com/jwcen/miniblog/internal/pkg/errno"
	"github.com/jwcen/miniblog/internal/pkg/log"
	"github.com/onexstack/onexstack/pkg/core"
)

type Authorizer interface {
	Authorize(sub, obj, act string) (bool, error)
}

func AuthzMiddleware(authorizer Authorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		sub := contextx.UserID(c)
		obj := c.Request.URL.Path
		act := c.Request.Method

		log.Debugw("Build authorize context", "subject", sub, "object", obj, "action", act)

		if allowed, err := authorizer.Authorize(sub, obj, act); err != nil || !allowed {
			core.WriteResponse(c, nil, errno.ErrPermissionDenied.WithMessage(
				"access denied: subject=%s, object=%s, action=%s, reason=%v",
				sub,
				obj,
				act,
				err,
			))
			c.Abort()
			return
		}

		c.Next()
	}
}
