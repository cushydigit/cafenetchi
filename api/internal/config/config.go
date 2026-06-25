package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	DBSSL  string

	RdsHost string
	RdsPort string
	RdsPass string

	ServerPort string
	JWTSecret  string

	KavenegarAPIKey string
	KavenegarSender string
}

func Load() *Config {
	// Load .env file in development
	_ = godotenv.Load()

	return &Config{
		DBHost: getEnv("DB_HOST", "postgres"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "cafenetchi-user"),
		DBPass: getEnv("DB_PASS", "cafenetchi-pass"),
		DBName: getEnv("DB_NAME", "cafenetchi-db"),
		DBSSL:  getEnv("DB_SSL", "disable"),

		RdsHost: getEnv("REDIS_HOST", "redis"),
		RdsPort: getEnv("REDIS_PORT", "6379"),
		RdsPass: getEnv("REDIS_PASS", ""),

		ServerPort:      getEnv("PORT", "8080"),
		JWTSecret:       getEnv("JWT_SECRET", ""),
		KavenegarAPIKey: getEnv("KAVENEGAR_API_KEY", ""),
		KavenegarSender: getEnv("KAVENEGAR_SENDER", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("⚠️  Environment variable %s not set, using default: %s", key, defaultValue)
	return defaultValue
}
