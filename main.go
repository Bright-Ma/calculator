package main

import (
	"calculator/internal/database"
	"calculator/internal/redis"
	"calculator/internal/router"
	"log"
)

func main() {
	// 初始化数据库连接
	if err := database.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化Redis连接
	redisClient := redis.NewRedis()

	// 初始化排行榜数据
	if err := redisClient.InitRankingData(); err != nil {
		log.Printf("初始化排行榜数据失败: %v", err)
	}

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
