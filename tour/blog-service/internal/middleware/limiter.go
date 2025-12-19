package middleware

import (
	"blog-service/pkg/app"
	"blog-service/pkg/errcode"
	"blog-service/pkg/limiter"

	"github.com/gin-gonic/gin"
)

func RateLimiter(l limiter.LimiterIface) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := l.Key(c)
		// 获取令牌桶
		if bucket, ok := l.GetBucket(key); ok {
			// 返回值为删除令牌数，如果删除成功，则返回1，否则返回0
			count := bucket.TakeAvailable(1)
			if count == 0 {
				// 返回错误信息
				response := app.NewResponse(c)
				response.ToErrorResponse(errcode.TooManyRequests)
				c.Abort() // 终止请求
				return
			}
		}
		// 如果令牌桶有令牌，则继续执行
		c.Next()
	}
}

