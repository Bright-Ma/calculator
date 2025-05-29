package model

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Password string `json:"-" gorm:"type:varchar(255);not null"`
	Role     string `json:"role" gorm:"type:varchar(20);not null"`
}
