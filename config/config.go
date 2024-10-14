package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"

	"github.com/joho/godotenv"
)

// Config holds the configuration for the database and the server

type Config struct {
	// Host is the database host
	DB_Host string `env:"DATABASE_HOST,required"`
	// Port is the database port
	DB_Port int `env:"DATABASE_PORT,required"`
	// User is the database user
	DB_User string `env:"DATABASE_USER,required"`
	// Password is the database password
	DB_Password string `env:"DATABASE_PASSWORD,required"`
	// Name is the database name
	DB_Name string `env:"DATABASE_NAME,required"`
	// RetryDuration is the duration to wait before retrying to connect to the database
	DB_RetryDuration string `env:"DATABASE_RETRY_DURATION,default=3s"`

	// Domain is the server domain
	HTTP_Domain string `env:"HTTP_DOMAIN,default=localhost"`
	// Port is the server port
	HTTP_Port string `env:"HTTP_PORT,default=8000"`
}


func New() (Config, error){
	var c Config
	godotenv.Load()
	if err := envconfig.Process(context.Background(), &c); err != nil {
		return Config{}, err
	}
	return c, nil
}