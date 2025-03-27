package model

import (
	"time"
)

type Message struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	CustomerID  uint64 `gorm:"index;not null" comment:"客户ID"`
	Message     string `gorm:"type:text;not null" json:"message" comment:"消息内容"`
	Sender      string `gorm:"size:32;not null" json:"sender" comment:"发送者标识"`
	MessageType int    `gorm:"not null" json:"message_type" comment:"消息类型"`

	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

const (
	MessageTypeNormal = iota
	MessageTypeFeedback
)

// TableName 自定义表名
func (Message) TableName() string {
	return "messages"
}
