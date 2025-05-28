package routes

import (
	"calculator/internal/database"
	"calculator/internal/handlers"
	"calculator/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *database.Database) *gin.Engine {
	router := gin.Default()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)

	// Auth routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
	}

	// Protected routes example
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Example protected route
		protected.GET("/profile", func(c *gin.Context) {
			userID := c.MustGet("userID").(uint)
			username := c.MustGet("username").(string)
			role := c.MustGet("role").(string)

			c.JSON(200, gin.H{
				"userID":   userID,
				"username": username,
				"role":     role,
			})
		})

		// Teacher-only route example
		teacher := protected.Group("/teacher")
		teacher.Use(middleware.RoleMiddleware("teacher"))
		{
			teacher.GET("/dashboard", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Welcome to teacher dashboard"})
			})
		}

		// Student-only route example
		student := protected.Group("/student")
		student.Use(middleware.RoleMiddleware("student"))
		{
			student.GET("/dashboard", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Welcome to student dashboard"})
			})
		}
	}

	return router
}
