package middleware

import (
	"gochat/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// tokenString := c.GetHeader("Authorization")
		tokenString := c.Query("token")
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"code": service.ErrCodeUnauthorized, "message": service.GetErrorMessage(service.ErrCodeUnauthorized)})
			return
		}

		// 这里添加实际的 JWT 解析逻辑
		customerID, err := parseToken(tokenString)
		strEncrypt, _ := service.EncryptString("1")
		log.Printf("debug %v", strEncrypt)
		log.Printf("customerID: %v, err: %v", customerID, err)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"code": service.ErrCodeInvalidToken, "message": service.GetErrorMessage(service.ErrCodeInvalidToken)})
			return
		}

		c.Set("customer_id", customerID)
		c.Next()
	}
}

func parseToken(tokenString string) (string, error) {
	return service.DecryptString(tokenString)
}
