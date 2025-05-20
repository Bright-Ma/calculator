package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"calculator/internal/calculator"
)

// Server 表示Web服务器
type Server struct {
	factory   *calculator.OperationFactory
	config    *calculator.Config
	generator *calculator.ProblemGenerator
}

// ProblemRequest 获取题目请求
type ProblemRequest struct {
	Difficulty string   `json:"difficulty"`
	Operations []string `json:"operations,omitempty"`
}

// ProblemResponse 题目响应
type ProblemResponse struct {
	Problem *calculator.Problem `json:"problem"`
	Error   string              `json:"error,omitempty"`
}

// AnswerRequest 提交答案请求
type AnswerRequest struct {
	ProblemID string  `json:"problemId"`
	Answer    float64 `json:"answer"`
	TimeSpent int     `json:"timeSpent"` // 答题用时(秒)
}

// AnswerResponse 答案响应
type AnswerResponse struct {
	Correct       bool    `json:"correct"`
	CorrectAnswer float64 `json:"correctAnswer,omitempty"`
	TimeSpent     int     `json:"timeSpent"`
	NeedRest      bool    `json:"needRest"`
	Error         string  `json:"error,omitempty"`
}

// ProblemSession 保存当前题目会话
type ProblemSession struct {
	Problem   *calculator.Problem
	StartTime time.Time
}

// 使用map存储当前活动的题目会话
var activeSessions = make(map[string]*ProblemSession)

// NewServer 创建新的Web服务器
func NewServer(factory *calculator.OperationFactory, config *calculator.Config) *Server {
	return &Server{
		factory:   factory,
		config:    config,
		generator: calculator.NewProblemGenerator(factory),
	}
}

// Start 启动Web服务器
func (s *Server) Start(addr string) error {
	// 静态文件服务
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// API路由
	http.HandleFunc("/api/problem/new", s.handleNewProblem)
	http.HandleFunc("/api/problem/answer", s.handleCheckAnswer)

	// Web界面路由
	http.HandleFunc("/", s.handleIndex)

	log.Printf("口算练习服务器启动在 %s\n", addr)
	return http.ListenAndServe(addr, nil)
}

// handleNewProblem 处理获取新题目的请求
func (s *Server) handleNewProblem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req ProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("解析请求失败: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("收到新题目请求: 难度=%s, 运算类型=%v", req.Difficulty, req.Operations)

	// 生成新题目
	problem, err := s.generator.GenerateProblem(req.Difficulty, req.Operations, 0)
	if err != nil {
		log.Printf("生成题目失败: %v", err)
		json.NewEncoder(w).Encode(ProblemResponse{Error: err.Error()})
		return
	}

	// 保存题目会话
	activeSessions[problem.ID] = &ProblemSession{
		Problem:   problem,
		StartTime: time.Now(),
	}

	// 清理超过30分钟的会话
	go s.cleanupOldSessions()

	json.NewEncoder(w).Encode(ProblemResponse{Problem: problem})
}

// handleCheckAnswer 处理检查答案的请求
func (s *Server) handleCheckAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req AnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("解析答案请求失败: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("收到答案检查请求: 题目ID=%s, 答案=%f, 用时=%d秒", req.ProblemID, req.Answer, req.TimeSpent)

	// 获取题目会话
	session, exists := activeSessions[req.ProblemID]
	if !exists {
		log.Printf("题目会话不存在: %s", req.ProblemID)
		json.NewEncoder(w).Encode(AnswerResponse{Error: "题目已过期"})
		return
	}

	// 检查答案
	correct := session.Problem.CheckAnswer(req.Answer)

	// 计算是否需要休息
	totalTime := int(time.Since(session.StartTime).Seconds())
	needRest := totalTime >= s.config.Problems.RestReminder

	// 删除已完成的会话
	delete(activeSessions, req.ProblemID)

	response := AnswerResponse{
		Correct:       correct,
		CorrectAnswer: session.Problem.Answer,
		TimeSpent:     req.TimeSpent,
		NeedRest:      needRest,
	}

	log.Printf("答案检查结果: 正确=%v, 用时=%d秒, 需要休息=%v", correct, req.TimeSpent, needRest)
	json.NewEncoder(w).Encode(response)
}

// cleanupOldSessions 清理过期的会话
func (s *Server) cleanupOldSessions() {
	threshold := time.Now().Add(-30 * time.Minute)
	for id, session := range activeSessions {
		if session.StartTime.Before(threshold) {
			delete(activeSessions, id)
		}
	}
}

// handleIndex 处理主页请求
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", "index.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Operations   []calculator.OperationConfig
		Difficulties []calculator.DifficultyConfig
		RestReminder int
	}{
		Operations:   s.config.GetEnabledOperations(),
		Difficulties: s.config.Problems.Difficulties,
		RestReminder: s.config.Problems.RestReminder,
	}

	log.Printf("渲染主页: %d个操作, %d个难度级别", len(data.Operations), len(data.Difficulties))
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
