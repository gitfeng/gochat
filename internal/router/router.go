// 假设使用Gin框架，以下是初始化router的示例代码
package router

import (
	"gochat/internal/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// 确保在路由注册顺序中，/ws 路由在静态文件路由之前
	// 如果有其他路由分组，请将ws路由放在最前面
	initChatRouter(r)
	initMessageRouter(r)

	// 添加健康检查路由
	r.GET("/healthcheck", handler.HealthCheckHandler)
}
