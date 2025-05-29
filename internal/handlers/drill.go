package handlers

import (
	"calculator/internal/database"
	"calculator/internal/drill"
	"calculator/internal/model"
	"calculator/internal/redis"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// Redis key 前缀
	questionKeyPrefix = "question:"
	dailyRankKey      = "rank:daily"
	weeklyRankKey     = "rank:weekly"
	// 时间衰减因子（24小时）
	timeDecayFactor = 24 * time.Hour
)

// 包级别默认 handler 实例，供路由直接调用
var defaultDrillHandler = &DrillHandler{
	generator: drill.NewGenerator(),
	redis:     redis.NewRedis(),
}

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// 注册热度排行榜路由
	r.GET("/api/drill/rankings", GetHotRanking)
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

	// 生成一个唯一的题目ID，包含时间戳和用户ID
	userID := int64(c.GetUint("user_id"))
	timestamp := time.Now().UnixNano()
	// 将用户ID和时间戳组合成一个唯一的ID
	questionID := (userID * 1000000000000) + (timestamp % 1000000000000)

	// 将题目序列化为JSON
	questionJSON, err := json.Marshal(question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "题目序列化失败"})
		return
	}

	// 存储题目到Redis，设置30分钟过期时间
	ctx := context.Background()
	err = defaultDrillHandler.redis.Client.Set(ctx,
		fmt.Sprintf("%s%d", redis.QuestionKeyPrefix, questionID),
		questionJSON,
		30*time.Minute,
	).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存题目失败"})
		return
	}

	// 返回给前端的数据格式
	c.JSON(http.StatusOK, gin.H{
		"id":         questionID,          // 题目ID
		"question":   question.Expression, // 题目表达式
		"difficulty": difficultyStr,       // 难度
	})
}

// SubmitAnswer 提交答案
func SubmitAnswer(c *gin.Context) {
	var req struct {
		QuestionID int64  `json:"question_id" binding:"required"`
		Answer     *int   `json:"answer" binding:"required"`
		Question   string `json:"question" binding:"required"`
		Difficulty string `json:"difficulty" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 从Redis获取题目
	ctx := context.Background()
	questionJSON, err := defaultDrillHandler.redis.Client.Get(ctx, fmt.Sprintf("%s%d", redis.QuestionKeyPrefix, req.QuestionID)).Bytes()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "题目不存在或已过期"})
		return
	}

	// 反序列化题目
	var question drill.Question
	if err := json.Unmarshal(questionJSON, &question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "题目数据无效"})
		return
	}

	// 判断答案是否正确
	isCorrect := *req.Answer == question.Answer

	// 创建历史记录
	history := model.HistoryRecord{
		UserID:           c.GetUint("user_id"),
		QuestionID:       fmt.Sprintf("%d", req.QuestionID),
		Question_content: req.Question,
		UserAnswer:       *req.Answer,
		CorrectAnswer:    question.Answer,
		IsCorrect:        isCorrect,
		Difficulty:       req.Difficulty,
		TimeSpent:        0, // 暂时不记录用时
	}

	if err := database.DB.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存历史记录失败"})
		return
	}

	// 更新用户热度值
	if err := defaultDrillHandler.redis.UpdateUserHotScore(ctx, c.GetUint("user_id"), isCorrect); err != nil {
		// 热度更新失败不影响答题结果
		fmt.Printf("更新热度失败: %v\n", err)
	}

	// 返回结果
	message := "回答正确！"
	if !isCorrect {
		message = fmt.Sprintf("回答错误，正确答案是：%d", question.Answer)
	}

	c.JSON(http.StatusOK, gin.H{
		"correct": isCorrect,
		"message": message,
	})
}

// GetHotRanking 获取热度排行榜
func GetHotRanking(c *gin.Context) {
	// 获取排行榜类型（小时榜/日榜）
	rankType := c.DefaultQuery("type", "hourly")
	if rankType != "hourly" && rankType != "daily" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的排行榜类型"})
		return
	}

	// 获取排行榜数据
	rankings, err := defaultDrillHandler.redis.GetHotRanking(c.Request.Context(), rankType, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取排行榜失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rankings": rankings,
	})
}

type DrillHandler struct {
	generator *drill.Generator
	redis     *redis.Redis
}
