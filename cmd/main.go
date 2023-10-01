package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tapiaw38/auth-api/internal/config"
	"github.com/tapiaw38/auth-api/internal/router"
	"github.com/tapiaw38/auth-api/internal/server"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error load env file")
	}

	config := config.NewConfig()

	s, err := server.NewServer(&server.Config{
		GinMode:            config.GinMode,
		Port:               config.Port,
		JWTSecret:          config.JWTSecret,
		DatabaseURL:        config.DatabaseURL,
		Host:               config.Host,
		AWSRegion:          config.AWSRegion,
		AWSAccessKeyID:     config.AWSAccessKeyID,
		AWSSecretAccessKey: config.AWSSecretAccessKey,
		AWSBucket:          config.AWSBucket,
		RedisHost:          config.RedisHost,
		RedisPassword:      config.RedisPassword,
		RedisDB:            config.RedisDB,
		RedisExpires:       config.RedisExpires,
		GoogleClientID:     config.GoogleClientID,
		GoogleClientSecret: config.GoogleClientSecret,
		FrontendURL:        config.FrontendURL,
		EmailHost:          config.EmailHost,
		EmailPort:          config.EmailPort,
		EmailHostUser:      config.EmailHostUser,
		EmailHostPassword:  config.EmailHostPassword,
		RabbitMQHost:       config.RabbitMQHost,
		RabbitMQPort:       config.RabbitMQPort,
		RabbitMQUser:       config.RabbitMQUser,
		RabbitMQPassword:   config.RabbitMQPassword,
	})

	if err != nil {
		log.Fatal(err)
	}

	s.Serve(router.BinderRoutes)
}
