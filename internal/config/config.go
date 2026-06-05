package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env      string       `mapstructure:"env"`
	Server   ServerConfig `mapstructure:"server"`
	Database DBConfig     `mapstructure:"db"`
	Auth     AuthConfig   `mapstructure:"auth"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type AuthConfig struct {
	URL string `mapstructure:"url" env:"AUTH_SERVICE_URL"`
}

func Load() (*Config, error) {
	_ = godotenv.Load() // игнорируем ошибку, если .env нет

	cfg := &Config{
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
		Auth: AuthConfig{
			URL: getEnv("AUTH_SERVICE_URL", "http://auth-service:5001"),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
