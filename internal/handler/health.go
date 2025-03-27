package handler

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HealthCheckHandler(c *gin.Context) {
	// 获取数据库连接
	db, ok := c.Get("DB")
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "database connection missing",
		})
		return
	}

	// 执行数据库健康检查
	sqlDB, err := db.(*gorm.DB).DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "db connection failed",
			"error":  err.Error(),
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "db ping failed",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": "1.0.0",
	})
}