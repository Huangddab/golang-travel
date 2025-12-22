package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// 统一超时控制中间件

// 接收一个超时时间，返回一个 gin.HandlerFunc 类型
// 在gin.Context中设置一个超时时间，如果超时，则返回一个错误
// 如果超时，则返回一个错误
// 如果超时，则返回一个错误
func ContextTimeout(t time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

