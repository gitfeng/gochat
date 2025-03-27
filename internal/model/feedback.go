package model

import (
	"time"
)

type Feedback struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	CustomerID uint64 `gorm:"index"`
	Score      uint   // requested, completed
	Comment    string `gorm:"type:text"`
	Sentiment  int    `gorm:"type:tinyint;default:0" json:"sentiment"`

	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 自定义表名
func (Feedback) TableName() string {
	return "feedback"
}
