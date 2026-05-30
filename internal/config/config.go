package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env      string
	Server   ServerConfig
	Database DBConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func Load() (*Config, error) {
	_ = godotenv.Load() // игнорируем ошибку, если .env нет

	return &Config{
		Env: getEnv("ENV", "development"),
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "albumsdb"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
