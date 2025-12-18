package middleware

import (
	"blog-service/pkg/app"
	"blog-service/pkg/errcode"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// JWT 中间件 从header中获取token并解析验证
func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			token string
			ecode = errcode.Success
		)
		// 优先从查询参数中获取 token，其次从 Header 中获取
		if s, exit := ctx.GetQuery("token"); exit {
			token = s
		} else {
			token = ctx.GetHeader("token")
		}

		// token 为空则返回错误
		if token == "" {
			// 缺少 token 应被视为鉴权失败
			ecode = errcode.UnauthorizedAuthNotExist
		} else {
			// 非空则解析 token
			_, err := app.ParseToken(token)
			// 解析出错则返回错误信息
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				// Token 过期
				case jwt.ValidationErrorExpired:
					ecode = errcode.UnauthorizedTokenTimeout
				// Token 无效
				default:
					ecode = errcode.UnauthorizedTokenError
				}
			}
		}
		// 如果有错误则返回错误信息并中止请求
		if ecode != errcode.Success {
			// 创建响应实例
			responses := app.NewResponse(ctx)
			responses.ToErrorResponse(ecode)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
