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

	s, err := server.NewServer(&server.Config{
		GinMode:            GIN_MODE,
		Port:               PORT,
		JWTSecret:          JW_SECRET,
		DatabaseURL:        DATABASE_URL,
		AWSRegion:          AWS_REGION,
		AWSAccessKeyID:     AWS_ACCESS_KEY_ID,
		AWSSecretAccessKey: AWS_SECRET_ACCESS_KEY,
		AWSBucket:          AWS_BUCKET,
		RedisHost:          REDIS_HOST,
		RedisPassword:      REDIS_PASSWORD,
		RedisDB:            0,
		RedisExpires:       10,
		GoogleClientID:     GOOGLE_CLIENT_ID,
		GoogleClientSecret: GOOGLE_CLIENT_SECRET,
		FrontendURL:        FRONTEND_URL,
	})

	if err != nil {
		log.Fatal(err)
	}

	s.Serve(router.BinderRoutes)
}
