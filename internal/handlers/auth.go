package handlers

import (
	"calculator/internal/database"
	"calculator/internal/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Register 用户注册
func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("注册请求解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供完整的注册信息"})
		return
	}

	log.Printf("收到注册请求: username=%s, role=%s", input.Username, input.Role)

	// 检查用户名是否已存在
	var existingUser model.User
	if err := database.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		log.Printf("用户名已存在: %s", input.Username)
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建用户
	user := model.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("创建用户失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	log.Printf("用户注册成功: username=%s", input.Username)
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// Login 用户登录
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("登录请求解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供用户名和密码"})
		return
	}

	log.Printf("收到登录请求: username=%s", input.Username)

	// 查找用户
	var user model.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		log.Printf("用户不存在: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		log.Printf("密码错误: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成 JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		log.Printf("生成token失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	log.Printf("用户登录成功: username=%s", user.Username)
	c.JSON(http.StatusOK, gin.H{
		"token":    tokenString,
		"username": user.Username,
		"role":     user.Role,
	})
}

// Logout 用户登出
func Logout(c *gin.Context) {
	// 由于使用 JWT，服务器端不需要维护会话状态
	// 客户端只需要删除本地存储的 token 即可
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}
