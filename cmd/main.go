package main

import (
	"github.com/tapiaw38/auth-api/config"
	"log"

	"github.com/joho/godotenv"
	"github.com/tapiaw38/auth-api/internal/router"
	"github.com/tapiaw38/auth-api/internal/server"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error load env file")
	}

	conf := config.New()

	s, err := server.New(conf)

	if err != nil {
		log.Fatal(err)
	}

	s.Serve(router.BinderRoutes)
}
