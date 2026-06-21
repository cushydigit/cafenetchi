package config

import (
	"cafenetchi-api/internal/model"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	DBSSL  string

	ServerPort string
	JWTSecret  string
}

func Load() *Config {
	// Load .env file in development
	_ = godotenv.Load()

	return &Config{
		DBHost:     getEnv("DB_HOST", "postgres"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "cafenetchi-user"),
		DBPass:     getEnv("DB_PASS", "cafenetchi-pass"),
		DBName:     getEnv("DB_NAME", "cafenetchi-db"),
		DBSSL:      getEnv("DB_SSL", "disable"),
		ServerPort: getEnv("PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "secret"),
	}
}

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

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
