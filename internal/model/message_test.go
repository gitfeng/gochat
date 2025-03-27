package model_test

import (
	"testing"
	"time"

	"gochat/internal/model"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestMessageValidation(t *testing.T) {
	dsn := viper.GetString("database.mysql.dsn")
	dsn = "root:123qwe@tcp(127.0.0.1:3306)/gochat?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("无法连接到数据库: %v", err)
	}
	db.AutoMigrate(&model.Message{})

	t.Run("有效消息", func(t *testing.T) {
		msg := model.Message{
			CustomerID:  1,
			Message:     "Valid message",
			Sender:      "user",
			MessageType: model.MessageTypeNormal,
			CreatedAt:   time.Now(),
		}
		result := db.Create(&msg)
		assert.NoError(t, result.Error)
	})
}
