package router

import (
	"gochat/internal/handler"
	"gochat/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ... 现有 upgrader 定义保持不变 ...
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 在这里进行权限验证，返回true表示验证通过，允许连接
		// 你可以根据需要实现自己的验证逻辑
		return true
	},
}

func initChatRouter(r *gin.Engine) {

	// 修改路由处理函数
	r.GET("/ws", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
		handler.ServeWebSocket(c, upgrader) // 传递gin上下文
	})
}
