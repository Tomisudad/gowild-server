package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 请求日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		log.Printf("[GIN] %s | %3d | %13v | %15s | %-7s %s | %d bytes",
			time.Now().Format("2006-01-02 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			bodySize,
		)

		// 慢请求告警
		if latency > time.Second {
			log.Printf("[SLOW] %s %s took %v", method, path, latency)
		}
	}
}
