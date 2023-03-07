package server

import (
	"errors"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/auth-api/internal/cache"
	"github.com/tapiaw38/auth-api/internal/database"
	"github.com/tapiaw38/auth-api/internal/repository"
	"github.com/tapiaw38/auth-api/internal/sso"
	"github.com/tapiaw38/auth-api/internal/utils"
)

// Config is the server configuration
type Config struct {
	GinMode              string
	Port                 string
	JWTSecret            string
	DatabaseURL          string
	Host                 string
	AWSRegion            string
	AWSAccessKeyID       string
	AWSSecretAccessKey   string
	AWSBucket            string
	RedisHost            string
	RedisPassword        string
	RedisDB              int
	RedisExpires         time.Duration
	GoogleClientID       string
	GoogleClientSecret   string
	FrontendURL          string
	EmailHost            string
	EmailPort            string
	EmailHostUser        string
	EmailHostPassword    string
	MailgunDomain        string
	MailgunPrivateAPIKey string
}

// Server is the server interface
type Server interface {
	Config() *Config
	S3() *utils.S3Client
	Google() *sso.GoogleClient
	Mail() *utils.EmailSMTPConfig
	Mailgun() *utils.MailgunConfig
	Redis() *cache.RedisCache
}

// Broker is the server broker
type Broker struct {
	config  *Config
	engine  *gin.Engine
	s3      *utils.S3Client
	google  *sso.GoogleClient
	mail    *utils.EmailSMTPConfig
	mailgun *utils.MailgunConfig
	redis   *cache.RedisCache
}

// Config returns the server configuration
func (b *Broker) Config() *Config {
	return b.config
}

// S3 returns the s3 client
func (b *Broker) S3() *utils.S3Client {
	return b.s3
}

// Google returns the google client
func (b *Broker) Google() *sso.GoogleClient {
	return b.google
}

// Mail returns the mail client
func (b *Broker) Mail() *utils.EmailSMTPConfig {
	return b.mail
}

// Mailgun returns the mailgun client
func (b *Broker) Mailgun() *utils.MailgunConfig {
	return b.mailgun
}

// Redis returns the redis client
func (b *Broker) Redis() *cache.RedisCache {
	return b.redis
}

// NewServer creates a new server
func NewServer(config *Config) (*Broker, error) {
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
		mailgun: utils.NewMailgunConfig(&utils.MailgunConfig{
			Domain:        config.MailgunDomain,
			PrivateAPIKey: config.MailgunPrivateAPIKey,
		}),
		redis: cache.NewRedisCache(&cache.RedisCache{
			Host:     config.RedisHost,
			Password: config.RedisPassword,
			DB:       config.RedisDB,
			Expires:  config.RedisExpires,
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
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowCredentials = true
	config.AllowMethods = []string{"*"}
	config.AllowHeaders = []string{"*"}
	config.ExposeHeaders = []string{"*"}

	// Use the cors
	b.engine.Use(cors.New(config))
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
