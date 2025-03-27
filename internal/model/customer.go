package model

import (
	"time"
)

type Customer struct {
	ID           uint64    `gorm:"primaryKey;column:customer_id;autoIncrement:true" json:"customer_id"`
	CustomerName string    `gorm:"type:varchar(128);uniqueIndex;not null" json:"customer_name"`
	Password     string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 自定义表名
func (Customer) TableName() string {
	return "customers"
}
