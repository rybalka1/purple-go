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
	AppEnv string
}

// Load загружает переменные окружения из .env и формирует конфигурацию
func Load() *Config {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Файл .env не найден, используются системные переменные")
	}

	cfg := &Config{
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPass: getEnv("DB_PASS", "password"),
		DBName: getEnv("DB_NAME", "order_api"),
		AppEnv: getEnv("APP_ENV", "development"),
	}

	return cfg
}

// getEnv возвращает значение переменной или дефолт, если не задано
func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
