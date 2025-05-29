package drill

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AnswerRequest 提交答案的请求结构
type AnswerRequest struct {
	QuestionID string `json:"question_id"` // 题目ID
	Answer     int    `json:"answer"`      // 用户答案
}

// AnswerResponse 提交答案的响应结构
type AnswerResponse struct {
	Correct bool   `json:"correct"` // 是否正确
	Score   int    `json:"score"`   // 得分
	Message string `json:"message"` // 反馈信息
}

// RegisterHandlers 注册口算题相关路由
func RegisterHandlers(r *gin.Engine, generator *Generator) {
	// 获取题目
	r.GET("/api/questions", func(c *gin.Context) {
		// 从查询参数获取难度，默认为简单
		difficulty := c.DefaultQuery("difficulty", "1")

		var d Difficulty
		switch difficulty {
		case "1":
			d = Easy
		case "2":
			d = Medium
		case "3":
			d = Hard
		default:
			d = Easy
		}

		question := generator.Generate(d)
		// 添加CORS头
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		c.Header("Content-Type", "application/json")
		jsonData, err := json.Marshal(map[string]interface{}{
			"question":   question.Expression,
			"difficulty": question.Difficulty,
			"answer":     question.Answer,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode JSON"})
			return
		}
		c.Data(http.StatusOK, "application/json", jsonData)
	})

	// 提交答案
	r.POST("/api/answers", func(c *gin.Context) {
		var req AnswerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 这里简化处理，实际应该验证答案
		resp := AnswerResponse{
			Correct: true, // 假设总是正确
			Score:   10,   // 固定得分
			Message: "回答正确!",
		}
		c.JSON(http.StatusOK, resp)
	})
}
