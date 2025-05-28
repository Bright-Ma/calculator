package main

import (
	"calculator/math-drill/backend/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

// 全局数据库连接
var db *sql.DB

// 配置文件结构
type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	DB struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
}

// 加载配置文件
func loadConfig(path string) (*Config, error) {
	config := &Config{}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// 初始化数据库连接
func initDB(config *Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DB.User,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// 精简版CORS中间件
func corsMiddleware(next routeHandler) routeHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// 认证中间件
func authMiddleware(next routeHandler) routeHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头获取token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No authorization header", http.StatusUnauthorized)
			return
		}

		// 验证token格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// 验证token有效性
		userID, valid := handlers.ValidateToken(token)
		if !valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 将用户ID添加到请求上下文
		r.Header.Set("X-User-ID", userID)

		// 继续处理请求
		next(w, r)
	}
}

// 路由处理函数类型
type routeHandler func(http.ResponseWriter, *http.Request)

// 路由配置结构体
type route struct {
	path      string
	handler   routeHandler
	needsAuth bool
	methods   []string
}

func main() {
	// 定义路由配置
	routes := []route{
		{
			path:      "/api/auth/register",
			handler:   handlers.Register,
			needsAuth: false,
			methods:   []string{"POST", "OPTIONS"},
		},
		{
			path:      "/api/auth/login",
			handler:   handlers.Login,
			needsAuth: false,
			methods:   []string{"POST", "OPTIONS"},
		},
		{
			path:      "/api/exercises",
			handler:   handlers.GetExercise,
			needsAuth: true,
			methods:   []string{"GET", "OPTIONS"},
		},
		{
			path:      "/api/records",
			handler:   handlers.SubmitAnswer,
			needsAuth: true,
			methods:   []string{"POST", "OPTIONS"},
		},
	}

	// 注册路由
	for _, route := range routes {
		handler := route.handler

		// 添加CORS中间件
		handler = corsMiddleware(handler)

		// 添加方法验证
		handler = methodMiddleware(handler, route.methods)

		// 如果需要认证，添加认证中间件
		if route.needsAuth {
			handler = authMiddleware(handler)
		}

		// 注册最终的处理函数
		http.HandleFunc(route.path, handler)
	}

	// 启动服务器
	// 假设使用 math-drill/config/config.yaml 中的端口配置
	port := 8080
	log.Printf("Config port value: %v (type: %T)", port, port)
	log.Printf("Server starting on :%d", port)
	log.Printf("Available routes:")
	for _, route := range routes {
		log.Printf("  %s [%s]", route.path, strings.Join(route.methods, ", "))
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// 方法验证中间件
func methodMiddleware(next routeHandler, allowedMethods []string) routeHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		// OPTIONS请求总是允许，用于CORS预检
		if r.Method == "OPTIONS" {
			next(w, r)
			return
		}

		// 检查请求方法是否在允许列表中
		methodAllowed := false
		for _, method := range allowedMethods {
			if r.Method == method {
				methodAllowed = true
				break
			}
		}

		if !methodAllowed {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		next(w, r)
	}
}
