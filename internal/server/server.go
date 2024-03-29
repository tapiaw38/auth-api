package server

import (
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api/config"
	"github.com/tapiaw38/auth-api/internal/cache"
	"github.com/tapiaw38/auth-api/internal/database"
	"github.com/tapiaw38/auth-api/internal/rabbitmq"
	"github.com/tapiaw38/auth-api/internal/repository"
	"github.com/tapiaw38/auth-api/internal/sso"
	"github.com/tapiaw38/auth-api/internal/utils"
	"log"
)

// Server is the server interface
type Server interface {
	Config() *config.Config
	S3() *utils.S3Client
	Google() *sso.GoogleClient
	Mail() *utils.EmailSMTPConfig
	Redis() *cache.RedisCache
	Rabbit() *rabbitmq.RabbitMQConfig
}

// Broker is the server broker
type Broker struct {
	config *config.Config
	engine *gin.Engine
	s3     *utils.S3Client
	google *sso.GoogleClient
	mail   *utils.EmailSMTPConfig
	redis  *cache.RedisCache
	rabbit *rabbitmq.RabbitMQConfig
}

// Config returns the server configuration
func (b *Broker) Config() *config.Config {
	return b.config
}

// S3 returns the s3 client
func (b *Broker) S3() *utils.S3Client {
	return b.s3
}

// Google returns the Google client
func (b *Broker) Google() *sso.GoogleClient {
	return b.google
}

// Mail returns the mail client
func (b *Broker) Mail() *utils.EmailSMTPConfig {
	return b.mail
}

// Redis returns the redis client
func (b *Broker) Redis() *cache.RedisCache {
	return b.redis
}

// Rabbit returns the rabbit client
func (b *Broker) Rabbit() *rabbitmq.RabbitMQConfig {
	return b.rabbit
}

// NewServer creates a new server
func New(config *config.Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("jwt secret is required")
	}

	if config.DatabaseURL == "" {
		return nil, errors.New("database url is required")
	}

	broker := &Broker{
		config: config,
		engine: gin.Default(),
		s3: utils.NewSession(&utils.S3Config{
			AWSRegion:          config.AWSRegion,
			AWSAccessKeyID:     config.AWSAccessKeyID,
			AWSSecretAccessKey: config.AWSSecretAccessKey,
			AWSBucket:          config.AWSBucket,
		}),
		google: sso.NewGoogleClient(&sso.GoogleClient{
			ClientID:     config.GoogleClientID,
			ClientSecret: config.GoogleClientSecret,
			FrontendURL:  config.FrontendURL,
		}),
		mail: utils.NewEmailSMTPConfig(&utils.EmailSMTPConfig{
			Host:         config.EmailHost,
			Port:         config.EmailPort,
			HostUser:     config.EmailHostUser,
			HostPassword: config.EmailHostPassword,
		}),
		redis: cache.NewRedisCache(&cache.RedisCache{
			Host:     config.RedisHost,
			Password: config.RedisPassword,
			DB:       config.RedisDB,
			Expires:  config.RedisExpires,
		}),
		rabbit: rabbitmq.NewRabbitMQConfig(&rabbitmq.RabbitMQConfig{
			Host:     config.RabbitMQHost,
			Port:     config.RabbitMQPort,
			User:     config.RabbitMQUser,
			Password: config.RabbitMQPassword,
		}),
	}

	return broker, nil
}

// Serve starts the server
func (b *Broker) Serve(binder func(s Server, e *gin.Engine)) {

	// Set gin mode
	if b.config.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
		b.config.Host = "http://localhost:" + b.config.Port
	}

	// Connect to RabbitMQ
	conn := b.rabbit.Connection()
	defer conn.Close()

	// Consumer for sending emails
	go func() {
		err := conn.ConsumeEmailMessage(b.mail.SendEmail)
		if err != nil {
			log.Fatalf("Failed to consume messages: %s", err)
		}
	}()

	// Create a new repository
	rep, err := database.NewPostresRepository(b.config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	err = rep.Makemigration(b.config.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// Ensure the base roles
	err = rep.EnsureRole()
	if err != nil {
		log.Println(err)
	}

	// Set the repository
	repository.SetRepository(rep)

	// Set the router as the default one shipped with Gin
	b.engine = gin.Default()

	// Set the cors
	conf := cors.DefaultConfig()
	conf.AllowOrigins = []string{"*"}
	conf.AllowCredentials = true
	conf.AllowMethods = []string{"*"}
	conf.AllowHeaders = []string{"*"}
	conf.ExposeHeaders = []string{"*"}

	// Use the cors
	b.engine.Use(cors.New(conf))
	//Use the recovery middleware
	b.engine.Use(gin.Recovery())
	// Use the logger middleware
	b.engine.Use(gin.Logger())

	// Set the router
	binder(b, b.engine)

	// Start and run the server
	err = b.engine.Run(":" + b.config.Port)
	if err != nil {
		log.Fatal(err)
	}

}
