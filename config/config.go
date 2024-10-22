package config

import (
	"os"
	"strconv"
	"time"
)

// Config struct
type Config struct {
	GinMode            string
	Port               string
	JWTSecret          string
	DatabaseURL        string
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSBucket          string
	RedisHost          string
	RedisPassword      string
	RedisDB            int
	RedisExpires       time.Duration
	GoogleClientID     string
	GoogleClientSecret string
	FrontendURL        string
	EmailHost          string
	EmailPort          string
	EmailHostUser      string
	EmailHostPassword  string
	RabbitMQHost       string
	RabbitMQPort       string
	RabbitMQUser       string
	RabbitMQPassword   string
	Host               string
	Domain             string
}

func New() *Config {
	return &Config{
		GinMode:            getEnv("GIN_MODE", "debug"),
		Port:               getEnv("PORT", "8080"),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		AWSRegion:          getEnv("AWS_REGION", ""),
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSBucket:          getEnv("AWS_BUCKET", ""),
		RedisHost:          getEnv("REDIS_HOST", ""),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:            getEnvAsInt("REDIS_DB", 0),
		RedisExpires:       getEnvAsTimeDuration("REDIS_EXPIRES", 10),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		FrontendURL:        getEnv("FRONTEND_URL", ""),
		EmailHost:          getEnv("EMAIL_HOST", ""),
		EmailPort:          getEnv("EMAIL_PORT", ""),
		EmailHostUser:      getEnv("EMAIL_HOST_USER", ""),
		EmailHostPassword:  getEnv("EMAIL_HOST_PASSWORD", ""),
		RabbitMQHost:       getEnv("RABBITMQ_HOST", ""),
		RabbitMQPort:       getEnv("RABBITMQ_PORT", ""),
		RabbitMQUser:       getEnv("RABBITMQ_USER", ""),
		RabbitMQPassword:   getEnv("RABBITMQ_PASSWORD", ""),
		Host:               getEnv("HOST", ""),
		Domain:             getEnv("DOMAIN", "localhost:8080"),
	}
}

// getEnv func
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvAsInt func
func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

// getEnvAsTimeDuration func
func getEnvAsTimeDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return time.Duration(i)
		}
	}
	return fallback
}
