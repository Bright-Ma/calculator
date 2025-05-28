package main

import (
	"calculator/internal/config"
	"calculator/internal/database"
	"calculator/internal/drill"
	"calculator/internal/routes"
	"log"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	db, err := database.New(cfg.DBConnectionString)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	// 初始化题目生成器
	generator := drill.NewGenerator()

	// 注册认证路由
	authRouter := routes.SetupRoutes(db)
	authRouter.Static("/static", "./frontend")

	// 注册题目路由
	drill.RegisterHandlers(authRouter, generator)

	// 启动服务器
	log.Println("口算题服务启动，监听 :8080")
	if err := authRouter.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
