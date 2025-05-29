package database

import (
	"calculator/internal/model"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	dsn := getEnv("DB_CONNECTION_STRING", "root:123456@tcp(127.0.0.1:3306)/calculator?charset=utf8mb4&parseTime=True&loc=Local")

	// 连接数据库
	var dbErr error
	DB, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		return fmt.Errorf("连接数据库失败: %v", dbErr)
	}

	// 自动迁移数据库表
	if migrateErr := DB.AutoMigrate(&model.User{}, &model.Session{}, &model.HistoryRecord{}); migrateErr != nil {
		return fmt.Errorf("数据库迁移失败: %v", migrateErr)
	}

	log.Println("数据库连接成功")
	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

type Database struct {
	db *gorm.DB
}

func New(connectionString string) (*Database, error) {
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(&model.User{}, &model.Session{}, &model.HistoryRecord{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate models: %v", err)
	}

	log.Println("Successfully connected to database using GORM")
	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Database) CreateUser(user *model.User) error {
	return d.db.Create(user).Error
}

func (d *Database) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := d.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (d *Database) CreateSession(session *model.Session) error {
	return d.db.Create(session).Error
}

func (d *Database) GetSessionByToken(token string) (*model.Session, error) {
	var session model.Session
	err := d.db.Where("token = ?", token).First(&session).Error
	return &session, err
}

func (d *Database) DeleteSession(token string) error {
	return d.db.Where("token = ?", token).Delete(&model.Session{}).Error
}
