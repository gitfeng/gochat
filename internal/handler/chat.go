package handler

import (
	"fmt"
	"gochat/internal/model"
	"gochat/internal/service"
	"gochat/internal/service/chatbot"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// 检测反馈关键词的辅助函数
// 新增配置加载部分
var (
	feedbackKeywords []string
	keywordsOnce     sync.Once
)

func init() {
	keywordsOnce.Do(func() {
		viper.SetConfigName("feedback_keywords")
		viper.AddConfigPath("./config")
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("无法读取关键词配置，使用默认值: %v", err)
			feedbackKeywords = []string{"feedback", "review"}
			return
		}
		feedbackKeywords = viper.GetStringSlice("triggers")
	})
}

// 在 ServeWebSocket 函数开头添加授权验证
func ServeWebSocket(c *gin.Context, upgrader websocket.Upgrader) {

	validCustomerID, err := validateSession(c)
	if err != nil {
		log.Println(err)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	// 设置读写超时（单位：秒）

	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	db := c.MustGet("DB").(*gorm.DB)

	chatbotEngine := chatbot.NewChatBotEngine(db)

	for {
		// 读取客户端消息
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		now := time.Now().Local()
		// 存储聊天记录
		msg := string(p)
		// 修改消息存储部分
		message := model.Message{
			CustomerID:  validCustomerID, // 替换原有硬编码 0
			Message:     msg,
			Sender:      "user", // 假设用户发送的消息为 "user"
			CreatedAt:   now,
			MessageType: model.MessageTypeNormal,
		}

		// 检测反馈关键词
		if containsFeedbackKeywords(msg) {
			message.MessageType = model.MessageTypeFeedback
			if result := db.Create(&message); result.Error != nil {
				log.Printf("Failed to save chat: %v", result.Error)
			}
			// 发送反馈提示
			feedbackPrompt := "Tanks for your feedback, we will deal in 3 days."
			if err := conn.WriteMessage(messageType, []byte(feedbackPrompt)); err != nil {
				log.Println(err)
				return
			}

			// 机器人消息也关联客户ID
			feedbackResponse := model.Message{
				CustomerID:  validCustomerID, // 替换原有硬编码 0
				Message:     feedbackPrompt,
				Sender:      "robot", // 假设用户发送的消息为 "user"
				MessageType: model.MessageTypeNormal,
				CreatedAt:   now,
			}
			if result := db.Create(&feedbackResponse); result.Error != nil {
				log.Printf("Failed to save chat: %v", result.Error)
			}

			// 机器人消息也关联客户ID
			feedback := model.Feedback{
				CustomerID: validCustomerID, // 替换原有硬编码 0
				Comment:    msg,
				CreatedAt:  now,
			}
			if result := db.Create(&feedback); result.Error != nil {
				log.Printf("Failed to save feedback: %v", result.Error)
			}
		} else {
			// 处理业务逻辑
			response := chatbotEngine.ProcessMessage(
				strconv.FormatUint(validCustomerID, 10),
				string(msg),
			)

			// 发送响应
			message := model.Message{
				CustomerID:  validCustomerID, // 替换原有硬编码 0
				Message:     response,
				Sender:      "robot", // 假设用户发送的消息为 "user"
				CreatedAt:   now,
				MessageType: model.MessageTypeNormal,
			}

			if result := db.Create(&message); result.Error != nil {
				log.Printf("Failed to save chat: %v", result.Error)
			}

			err = conn.WriteMessage(websocket.TextMessage, []byte(response))

			if err != nil {
				log.Println(err)
				return
			}
		}
	}

}

func validateSession(c *gin.Context) (uint64, error) {
	// 新增授权验证
	customerID, ok := c.Get("customer_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": service.ErrCodeUnauthorized, "message": service.GetErrorMessage(service.ErrCodeUnauthorized)})
		return 0, fmt.Errorf("未授权访问")
	}

	// 类型转换验证
	validCustomerID, err := strconv.ParseUint(fmt.Sprintf("%v", customerID), 10, 64)
	if err != nil {
		// 类型转换验证失败，返回无效用户ID
		c.JSON(http.StatusBadRequest, gin.H{"code": service.ErrCodeUserNotFound, "message": service.GetErrorMessage(service.ErrCodeUserNotFound)})
		return 0, fmt.Errorf("无效用户ID")
	}

	// 获取数据库连接
	db := c.MustGet("DB").(*gorm.DB)

	// 检查customer表中是否存在该用户ID
	var customerCount int64
	if err := db.Model(&model.Customer{}).Where("id = ?", validCustomerID).Count(&customerCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": service.ErrCodeInternalServer, "message": service.GetErrorMessage(service.ErrCodeInternalServer)})
		return validCustomerID, fmt.Errorf("数据库错误")
	}

	if customerCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": service.ErrCodeUserNotFound, "message": service.GetErrorMessage(service.ErrCodeUserNotFound)})
		return validCustomerID, fmt.Errorf("无效用户ID")
	}

	return validCustomerID, nil
}

// 修改后的关键词检测函数
func containsFeedbackKeywords(message string) bool {
	lowerMsg := strings.ToLower(message)

	for _, keyword := range feedbackKeywords {
		if strings.Contains(lowerMsg, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}
