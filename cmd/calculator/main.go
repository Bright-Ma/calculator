package main

import (
	"log"
	"path/filepath"

	"calculator/internal/calculator"
	"calculator/internal/server"
)

func main() {
	// 加载XML配置
	config, err := calculator.LoadConfig(filepath.Join("config", "calculator.xml"))
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建操作工厂
	factory := calculator.NewOperationFactory()

	// 创建Web服务器
	srv := server.NewServer(factory, config)

	// 启动服务器
	log.Println("启动计算器Web服务...")
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
