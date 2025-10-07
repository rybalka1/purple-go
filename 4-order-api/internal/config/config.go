package config

import (
	"log"
	"os"
)

type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func Load() *Config {
	cfg := &Config{
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPass: getEnv("DB_PASS", "password"),
		DBName: getEnv("DB_NAME", "order_api"),
	}

	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func MustEnv() {
	for _, k := range []string{"DB_HOST", "DB_USER", "DB_PASS"} {
		if os.Getenv(k) == "" {
			log.Fatalf("Environment variable %s is required", k)
		}
	}
}
