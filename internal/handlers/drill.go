package handlers

import (
	"calculator/internal/database"
	"calculator/internal/drill"
	"calculator/internal/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 包级别默认 handler 实例，供路由直接调用
var defaultDrillHandler = &DrillHandler{
	generator: drill.NewGenerator(),
}

// GetQuestion 获取一道新题目
func GetQuestion(c *gin.Context) {
	difficultyStr := c.DefaultQuery("difficulty", "easy")
	var difficulty drill.Difficulty

	switch difficultyStr {
	case "easy":
		difficulty = drill.Easy
	case "medium":
		difficulty = drill.Medium
	case "hard":
		difficulty = drill.Hard
	default:
		difficulty = drill.Easy
	}

	question := defaultDrillHandler.generator.Generate(difficulty)

	c.JSON(http.StatusOK, gin.H{
		"id":         question.Answer, // 使用答案作为ID
		"question":   question.Expression,
		"answer":     question.Answer,
		"difficulty": difficultyStr,
	})
}

// SubmitAnswer 提交答案
func SubmitAnswer(c *gin.Context) {
	var req struct {
		QuestionID int    `json:"question_id" binding:"required"`
		Answer     int    `json:"answer" binding:"required"`
		Question   string `json:"question" binding:"required"`
		Difficulty string `json:"difficulty" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 判断答案是否正确
	isCorrect := req.Answer == req.QuestionID

	// 创建历史记录
	history := model.HistoryRecord{
		UserID:           c.GetUint("user_id"),
		QuestionID:       fmt.Sprintf("%d", req.QuestionID),
		Question_content: req.Question,
		UserAnswer:       req.Answer,
		CorrectAnswer:    req.QuestionID, // QuestionID 就是正确答案
		IsCorrect:        isCorrect,
		Difficulty:       req.Difficulty,
		TimeSpent:        0, // 暂时不记录用时
	}

	if err := database.DB.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存历史记录失败"})
		return
	}

	// 返回结果
	message := "回答正确！"
	if !isCorrect {
		message = fmt.Sprintf("回答错误，正确答案是：%d", req.QuestionID)
	}

	c.JSON(http.StatusOK, gin.H{
		"correct": isCorrect,
		"message": message,
	})
}

type DrillHandler struct {
	generator *drill.Generator
}
