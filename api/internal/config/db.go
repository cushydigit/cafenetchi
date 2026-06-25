package config

import (
	"cafenetchi-api/internal/model"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort, cfg.DBSSL,
	)

	db, err := connectDB(dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB = db
	fmt.Println("✅ Database connected successfully")
	return db
}

func connectDB(dsn string) (*gorm.DB, error) {
	var err error

	for i := 0; i < 5; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			return db, nil
		} else {
			log.Printf("Failed to connect to database, retrying in 5 seconds... (attempt %d/5)", i+1)
			time.Sleep(5 * time.Second)
		}

	}
	return nil, fmt.Errorf("failed to connect to database after 5 attempts: %v", err)
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.UserProfile{},
		&model.Conversation{},
		&model.Message{},
		// Add new models here later
	)
}
