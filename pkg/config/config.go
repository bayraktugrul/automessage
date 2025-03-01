package config

import (
	"os"
	"strconv"
	"time"
)

type RootConfig struct {
	Database DatabaseConfig
	Redis    RedisConfig
	App      AppConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type AppConfig struct {
	Port          string
	WebhookURL    string
	MessageConfig MessageConfig
}

type MessageConfig struct {
	InitialBatchSize  int
	PeriodicBatchSize int
	IntervalSecond    time.Duration
}

func LoadConfig() RootConfig {
	return RootConfig{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "automsg"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		App: AppConfig{
			Port:       getEnv("PORT", "8080"),
			WebhookURL: getEnv("WEBHOOK_URL", ""),
			MessageConfig: MessageConfig{
				InitialBatchSize:  getEnvAsInt("MESSAGE_INITIAL_BATCH_SIZE", 10),
				PeriodicBatchSize: getEnvAsInt("MESSAGE_PERIODIC_BATCH_SIZE", 2),
				IntervalSecond:    time.Duration(getEnvAsInt("MESSAGE_INTERVAL_SECONDS", 120)) * time.Second,
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
