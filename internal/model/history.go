package model

import "time"

// HistoryRecord 历史记录模型
type HistoryRecord struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	UserID           uint      `json:"user_id" gorm:"not null"`
	QuestionID       string    `json:"question_id" gorm:"not null"`
	Question_content string    `json:"question_content" gorm:"not null"`
	UserAnswer       int       `json:"user_answer" gorm:"not null"`
	CorrectAnswer    int       `json:"correct_answer" gorm:"not null"`
	IsCorrect        bool      `json:"is_correct" gorm:"not null"`
	Difficulty       string    `json:"difficulty" gorm:"not null"`
	TimeSpent        float64   `json:"time_spent" gorm:"not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"not null"`
}
