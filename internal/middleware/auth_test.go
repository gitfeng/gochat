package middleware_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gochat/internal/middleware"
	"gochat/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	router := gin.Default()
	router.Use(middleware.JWTAuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("无认证头", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	strEncrypt, _ := service.EncryptString("1")
	log.Printf("debug %v", strEncrypt)
	t.Run("有效认证头", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected?token="+strEncrypt, nil)

		// 需要设置mock验证逻辑
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
