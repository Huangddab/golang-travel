package middleware

import (
	"blog-service/global"
	"blog-service/pkg/logger"
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessLogWriter struct {
	gin.ResponseWriter               // 继承 gin.ResponseWriter
	body               *bytes.Buffer // 用于缓存响应 body
}

// Write 重写 Write 方法，同时将响应写入 body 缓存
func (w AccessLogWriter) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

// 初始化AccessLogWriter并在中间件中使用
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建 AccessLogWriter 实例，包装原有的 ResponseWriter
		bodyWriter := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyWriter

		// 继续处理请求
		beginTime := time.Now().Unix()
		c.Next()
		endTime := time.Now().Unix()

		fields := logger.Fields{
			"request":  c.Request.PostForm.Encode(), // 当前的请求参数
			"response": bodyWriter.body.String(),    // 当前的请求结果响应
		}

		global.Logger.WithFileds(fields).Infof("access log: method: %s, status_code: %d, begin_time: %d, end_time: %d",
			c.Request.Method,    // 当前的调用方法
			bodyWriter.Status(), // 当前的响应状态码
			beginTime,           // 请求开始时间
			endTime,             // 请求结束时间
		)
	}
}
