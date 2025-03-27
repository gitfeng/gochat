package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 新增追踪上下文键
const traceKey = "X-Trace-ID"

// TraceMiddleware 全链路追踪中间件
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 生成或获取追踪ID
		traceID := c.GetHeader(traceKey)
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// 设置上下文和响应头
		c.Set(traceKey, traceID)
		c.Writer.Header().Set(traceKey, traceID)

		// 记录请求开始日志
		log.Printf("[%s] %s %s [START]", traceID, c.Request.Method, c.Request.URL.Path)

		// 处理请求
		c.Next()

		// 记录请求完成日志
		latency := time.Since(start)
		log.Printf("[%s] %s %s [END] status=%d latency=%v",
			traceID,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
		)
	}
}

// GetTraceID 从上下文中获取追踪ID
func GetTraceID(c *gin.Context) string {
	if id, exists := c.Get(traceKey); exists {
		return id.(string)
	}
	return ""
}
