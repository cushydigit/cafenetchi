package config

import (
	"cafenetchi-api/internal/logger"
	"fmt"
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

func Load(l *logger.Logger) *Config {
	appEnv := getEnv("APP_ENV", "development", l)
	// Load .env file in development
	if appEnv == "development" {
		l.Warn(
			"app running on development mode",
			"message", "⚠️ app running on development mode")
		err := godotenv.Load("../configs/dev.env")
		if err != nil {
			l.Error(
				"error loading .env file",
				"message",
				err,
			)
		}
	}

	return &Config{
		AppEnv:    appEnv,
		Port:      getEnv("PORT", "8080", l),
		JWTSecret: getEnv("JWT_SECRET", "", l),

		DB: DBConfig{
			Host: getEnv("DB_HOST", "localhost", l),
			Port: getEnv("DB_PORT", "5432", l),
			User: getEnv("DB_USER", "", l),
			Pass: getEnv("DB_PASS", "", l),
			Name: getEnv("DB_NAME", "", l),
			SSL:  getEnv("DB_SSL", "disable", l),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost", l),
			Port: getEnv("REDIS_PORT", "6379", l),
			Pass: getEnv("REDIS_PASS", "", l),
		},

		KavenegarAPIKey: getEnv("KAVENEGAR_API_KEY", "", l),
		KavenegarSender: getEnv("KAVENEGAR_SENDER", "", l),
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

func getEnv(key, defaultValue string, l *logger.Logger) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	l.Warn(
		"Environment variable not set",
		"key",
		key,
	)
	return defaultValue
}
