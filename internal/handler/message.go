package handler

import (
	"gochat/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MessageResponse struct {
	CustomerID uint64 `json:"customer_id"`
	Message    string `json:"message"`
	Sender     string `json:"sender"`

	CreatedAt time.Time `json:"timestamp"`
}

// 新增分页响应结构体
type PaginationResponse struct {
	Data       []MessageResponse `json:"data"`
	Pagination PaginationMeta    `json:"pagination"`
}

type PaginationMeta struct {
	Total       int `json:"total"`
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalPages  int `json:"total_pages"`
}

func GetMessageList(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var messages []model.Message

	// 获取分页参数，默认每页10条，第一页
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// 计算偏移量
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		// 处理错误，例如设置默认值
		pageInt = 1
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 || limitInt > 100 {
		// 处理错误，例如设置默认值
		limitInt = 10
	}

	offset := (pageInt - 1) * limitInt

	// 处理查询参数
	query := db.Model(&model.Message{})

	// 支持按customer_id查询
	customerID := c.Query("customer_id")
	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	// 新增总记录数查询
	var total int64
	queryCount := db.Model(&model.Message{})
	if customerID != "" {
		queryCount = queryCount.Where("customer_id = ?", customerID)
	}
	queryCount.Count(&total) // 获取总记录数

	// 按id倒序排序
	query = query.Order("id DESC").Limit(limitInt).Offset(offset)

	if err := query.Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 分页元数据计算
	totalPages := (int(total) + limitInt - 1) / limitInt
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	// 修改返回数据结构部分
	var responseData []MessageResponse
	for _, msg := range messages {
		responseData = append(responseData, MessageResponse{
			CustomerID: msg.CustomerID,
			Message:    msg.Message,
			Sender:     msg.Sender,
			CreatedAt:  msg.CreatedAt, // 保持时间字段自动转换
		})
	}

	// 更新分页响应结构
	c.JSON(http.StatusOK, PaginationResponse{
		Data: responseData, // 使用转换后的数据
		Pagination: PaginationMeta{
			Total:       int(total),
			CurrentPage: pageInt,
			PageSize:    limitInt,
			TotalPages:  totalPages,
		},
	})
}
