package model

import (
	"time"

	"gorm.io/gorm"
)

// Session 会话模型
type Session struct {
	gorm.Model
	UserID    uint      `gorm:"not null" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;size:255;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
}

// TableName 指定表名
func (Session) TableName() string {
	return "sessions"
}
