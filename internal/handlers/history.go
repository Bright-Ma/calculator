package handlers

import (
	"calculator/internal/database"
	"calculator/internal/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetHistory 获取用户的历史记录
func GetHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	difficulty := c.Query("difficulty")
	date := c.Query("date")

	var records []model.HistoryRecord
	query := database.DB.Where("user_id = ?", userID)

	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}

	if date != "" {
		startTime, err := time.Parse("2006-01-02", date)
		if err == nil {
			endTime := startTime.Add(24 * time.Hour)
			query = query.Where("created_at BETWEEN ? AND ?", startTime, endTime)
		}
	}

	if err := query.Order("created_at DESC").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取历史记录失败"})
		return
	}

	c.JSON(http.StatusOK, records)
}

// GetStatistics 获取用户的练习统计信息
func GetStatistics(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 获取不同难度的题目数量
	var stats struct {
		TotalQuestions  int64   `json:"total_questions"`  // 去重后的总题数
		EasyQuestions   int64   `json:"easy_questions"`   // 去重后的简单题数
		MediumQuestions int64   `json:"medium_questions"` // 去重后的中等题数
		HardQuestions   int64   `json:"hard_questions"`   // 去重后的难题数
		TotalAttempts   int64   `json:"total_attempts"`   // 总答题次数
		CorrectAnswers  int64   `json:"correct_answers"`  // 正确答题次数
		Accuracy        float64 `json:"accuracy"`         // 正确率（正确次数/总次数）
	}

	// 获取去重后的总题数
	if err := database.DB.Model(&model.HistoryRecord{}).
		Where("user_id = ?", userID).
		Distinct("question_id").
		Count(&stats.TotalQuestions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}

	// 获取去重后的各难度题目数
	if err := database.DB.Model(&model.HistoryRecord{}).
		Where("user_id = ? AND difficulty = ?", userID, "easy").
		Distinct("question_id").
		Count(&stats.EasyQuestions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}

	if err := database.DB.Model(&model.HistoryRecord{}).
		Where("user_id = ? AND difficulty = ?", userID, "medium").
		Distinct("question_id").
		Count(&stats.MediumQuestions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}

	if err := database.DB.Model(&model.HistoryRecord{}).
		Where("user_id = ? AND difficulty = ?", userID, "hard").
		Distinct("question_id").
		Count(&stats.HardQuestions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}

	// 获取总答题次数
	if err := database.DB.Model(&model.HistoryRecord{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalAttempts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}

	// 获取正确答题次数
	if err := database.DB.Model(&model.HistoryRecord{}).
		Where("user_id = ? AND is_correct = true", userID).
		Count(&stats.CorrectAnswers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}

	// 计算正确率（使用总答题次数作为分母）
	if stats.TotalAttempts > 0 {
		stats.Accuracy = float64(stats.CorrectAnswers) / float64(stats.TotalAttempts) * 100
	}

	c.JSON(http.StatusOK, stats)
}

// AddHistory 添加历史记录
func AddHistory(c *gin.Context) {
	var record model.HistoryRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	record.UserID = c.GetUint("user_id")
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	if err := database.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存历史记录失败"})
		return
	}

	c.JSON(http.StatusOK, record)
}
