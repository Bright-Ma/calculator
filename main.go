package main

import (
	"calculator/internal/database"
	"calculator/internal/router"
	"log"
)

func main() {
	// 初始化数据库连接
	database.InitDB()

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
