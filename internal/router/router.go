package router

import (
	"calculator/internal/handlers"
	"calculator/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置所有路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 允许跨域
	r.Use(middleware.CORS())

	// 设置静态文件目录
	r.Static("/static", "./static")
	r.StaticFile("/", "./frontend/index.html")
	r.StaticFile("/index.html", "./frontend/index.html")
	r.StaticFile("/style.css", "./frontend/style.css")
	r.StaticFile("/app.js", "./frontend/app.js")
	r.StaticFile("/auth.js", "./frontend/auth.js")
	r.StaticFile("/history.html", "./frontend/history.html")
	r.StaticFile("/history.js", "./frontend/history.js")
	r.StaticFile("/history-detail.html", "./frontend/history-detail.html")
	r.StaticFile("/history-detail.js", "./frontend/history-detail.js")

	// API 路由组
	api := r.Group("/api")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)
		}

		// 题目相关路由
		drill := api.Group("/drill")
		drill.Use(middleware.AuthRequired())
		{
			drill.GET("/question", handlers.GetQuestion)
			drill.POST("/answer", handlers.SubmitAnswer)
			drill.GET("/rankings", handlers.GetHotRanking)
		}

		// 历史记录相关路由
		history := api.Group("/history")
		history.Use(middleware.AuthRequired())
		{
			history.GET("", handlers.GetHistory)
			history.GET("/stats", handlers.GetStatistics)
			history.POST("", handlers.AddHistory)
		}
	}

	// 添加通配符路由，支持前端路由
	r.NoRoute(func(c *gin.Context) {
		// 如果是 API 请求，返回 404
		if c.Request.URL.Path[:4] == "/api" {
			c.JSON(404, gin.H{"error": "API not found"})
			return
		}
		// 否则返回前端页面
		c.File("./frontend/index.html")
	})

	return r
}
