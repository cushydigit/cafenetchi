package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv    string
	Port      string
	JWTSecret string

	DB    DBConfig
	Redis RedisConfig

	KavenegarAPIKey string
	KavenegarSender string
}

type DBConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
	SSL  string
}

type RedisConfig struct {
	Host string
	Port string
	Pass string
}

func Load() *Config {
	appEnv := getEnv("APP_ENV", "development")
	// Load .env file in development
	if appEnv == "development" {
		log.Println("⚠️ app running on development mode")
		err := godotenv.Load("../configs/dev.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	return &Config{
		AppEnv:    appEnv,
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", ""),

		DB: DBConfig{
			Host: getEnv("DB_HOST", "localhost"),
			Port: getEnv("DB_PORT", "5432"),
			User: getEnv("DB_USER", ""),
			Pass: getEnv("DB_PASS", ""),
			Name: getEnv("DB_NAME", ""),
			SSL:  getEnv("DB_SSL", "disable"),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnv("REDIS_PORT", "6379"),
			Pass: getEnv("REDIS_PASS", ""),
		},

		KavenegarAPIKey: getEnv("KAVENEGAR_API_KEY", ""),
		KavenegarSender: getEnv("KAVENEGAR_SENDER", ""),
	}
}

func (c RedisConfig) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)

}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.Name,
	)

}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("⚠️  Environment variable %s not set, using default: %s", key, defaultValue)
	return defaultValue
}
