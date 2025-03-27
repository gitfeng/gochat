package chatbot_test

import (
	"gochat/internal/service/chatbot"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestEngine(t *testing.T) (*chatbot.ChatBotEngine, *gorm.DB) {

	// 配置viper测试路径
	viper.Reset()
	viper.SetConfigName("chatbot_rules")
	viper.AddConfigPath("d:/workspace/gochat/config") // 绝对路径
	//viper.AddConfigPath("../../config") // 相对路径

	// 内存数据库
	dsn := "root:123qwe@tcp(127.0.0.1:3306)/gochat?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("无法连接到数据库: %v", err)
	}

	// 模拟配置加载
	engine := chatbot.NewChatBotEngine(db)
	return engine, db
}

func TestGreetingIntent(t *testing.T) {
	engine, _ := setupTestEngine(t)

	// 首次对话
	resp := engine.ProcessMessage("1001", "hello")
	assert.Equal(t, "您好，我是${bot_name}，请问需要什么帮助？", resp)

	// 验证上下文状态
	ctx := engine.GetContext("1001")
	assert.Equal(t, "welcome", ctx.CurrentState)
}

func TestUnknownIntent(t *testing.T) {
	engine, _ := setupTestEngine(t)

	resp := engine.ProcessMessage("1002", "Random message")
	assert.Contains(t, resp, "抱歉，我还在学习中，暂时无法回答这个问题")

	ctx := engine.GetContext("1002")
	assert.Equal(t, "welcome", ctx.CurrentState) // 应保持初始状态
}
