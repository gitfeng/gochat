package router

import (
	"gochat/internal/handler"

	"github.com/gin-gonic/gin"
)

func initMessageRouter(r *gin.Engine) {
	// 路由分组
	api := r.Group("/message/")
	{
		api.GET("/list", handler.GetMessageList)
	}
}
