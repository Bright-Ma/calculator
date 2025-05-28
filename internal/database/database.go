package database

import (
	"calculator/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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
	err = db.AutoMigrate(&models.User{}, &models.Session{})
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

func (d *Database) CreateUser(user *models.User) error {
	return d.db.Create(user).Error
}

func (d *Database) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := d.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (d *Database) CreateSession(session *models.Session) error {
	return d.db.Create(session).Error
}

func (d *Database) GetSessionByToken(token string) (*models.Session, error) {
	var session models.Session
	err := d.db.Where("token = ?", token).First(&session).Error
	return &session, err
}

func (d *Database) DeleteSession(token string) error {
	return d.db.Where("token = ?", token).Delete(&models.Session{}).Error
}
