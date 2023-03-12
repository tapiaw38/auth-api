package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tapiaw38/auth-api/internal/router"
	"github.com/tapiaw38/auth-api/internal/server"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error load env file")
	}

	// APP
	GIN_MODE := os.Getenv("GIN_MODE")
	PORT := os.Getenv("PORT")
	JW_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	// AWS
	AWS_REGION := os.Getenv("AWS_REGION")
	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_BUCKET := os.Getenv("AWS_BUCKET")

	// REDIS
	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")

	// GOOGLE
	GOOGLE_CLIENT_ID := os.Getenv("GOOGLE_CLIENT_ID")
	GOOGLE_CLIENT_SECRET := os.Getenv("GOOGLE_CLIENT_SECRET")
	FRONTEND_URL := os.Getenv("FRONTEND_URL")

	// MAILGUN
	MAILGUN_DOMAIN := os.Getenv("MAILGUN_DOMAIN")
	MAILGUN_API_KEY := os.Getenv("MAILGUN_API_KEY")

	// STPM MAIL
	EMAIL_HOST := os.Getenv("EMAIL_HOST")
	EMAIL_PORT := os.Getenv("EMAIL_PORT")
	EMAIL_HOST_USER := os.Getenv("EMAIL_HOST_USER")
	EMAIL_HOST_PASSWORD := os.Getenv("EMAIL_HOST_PASSWORD")

	// RABBITMQ
	RABBITMQ_HOST := os.Getenv("RABBITMQ_HOST")
	RABBITMQ_PORT := os.Getenv("RABBITMQ_PORT")
	RABBITMQ_USER := os.Getenv("RABBITMQ_USER")
	RABBITMQ_PASSWORD := os.Getenv("RABBITMQ_PASSWORD")

	HOST := os.Getenv("HOST")

	s, err := server.NewServer(&server.Config{
		GinMode:              GIN_MODE,
		Port:                 PORT,
		JWTSecret:            JW_SECRET,
		DatabaseURL:          DATABASE_URL,
		Host:                 HOST,
		AWSRegion:            AWS_REGION,
		AWSAccessKeyID:       AWS_ACCESS_KEY_ID,
		AWSSecretAccessKey:   AWS_SECRET_ACCESS_KEY,
		AWSBucket:            AWS_BUCKET,
		RedisHost:            REDIS_HOST,
		RedisPassword:        REDIS_PASSWORD,
		RedisDB:              0,
		RedisExpires:         10,
		GoogleClientID:       GOOGLE_CLIENT_ID,
		GoogleClientSecret:   GOOGLE_CLIENT_SECRET,
		FrontendURL:          FRONTEND_URL,
		EmailHost:            EMAIL_HOST,
		EmailPort:            EMAIL_PORT,
		EmailHostUser:        EMAIL_HOST_USER,
		EmailHostPassword:    EMAIL_HOST_PASSWORD,
		MailgunDomain:        MAILGUN_DOMAIN,
		MailgunPrivateAPIKey: MAILGUN_API_KEY,
		RabbitMQHost:         RABBITMQ_HOST,
		RabbitMQPort:         RABBITMQ_PORT,
		RabbitMQUser:         RABBITMQ_USER,
		RabbitMQPassword:     RABBITMQ_PASSWORD,
	})

	if err != nil {
		log.Fatal(err)
	}

	s.Serve(router.BinderRoutes)
}
