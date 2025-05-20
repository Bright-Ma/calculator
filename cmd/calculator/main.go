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

	log.Printf("配置加载成功: %d个操作, %d个难度级别",
		len(config.Operations),
		len(config.Problems.Difficulties))

	// 创建操作工厂
	factory := calculator.NewOperationFactory()
	log.Printf("操作工厂初始化完成")

	// 创建Web服务器
	srv := server.NewServer(factory, config)
	log.Printf("Web服务器初始化完成")

	// 启动服务器
	log.Println("启动小学生口算练习系统...")
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
