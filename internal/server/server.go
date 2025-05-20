package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"calculator/internal/calculator"
)

// Server 表示Web服务器
type Server struct {
	factory *calculator.OperationFactory
	config  *calculator.Config
}

// CalculateRequest 计算请求结构
type CalculateRequest struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

// CalculateResponse 计算响应结构
type CalculateResponse struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

// NewServer 创建新的Web服务器
func NewServer(factory *calculator.OperationFactory, config *calculator.Config) *Server {
	return &Server{
		factory: factory,
		config:  config,
	}
}

// Start 启动Web服务器
func (s *Server) Start(addr string) error {
	// 静态文件服务
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// API路由
	http.HandleFunc("/api/calculate", s.handleCalculate)

	// Web界面路由
	http.HandleFunc("/", s.handleIndex)

	log.Printf("服务器启动在 %s\n", addr)
	return http.ListenAndServe(addr, nil)
}

// handleIndex 处理主页请求
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", "index.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Operations []calculator.OperationConfig
	}{
		Operations: s.config.GetEnabledOperations(),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleCalculate 处理计算请求
func (s *Server) handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	operation, err := s.factory.Create(req.Operation)
	if err != nil {
		json.NewEncoder(w).Encode(CalculateResponse{
			Error: err.Error(),
		})
		return
	}

	result := operation.Calculate(req.A, req.B)
	json.NewEncoder(w).Encode(CalculateResponse{
		Result: result,
	})
}
