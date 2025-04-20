package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Db   DbConfig
	Auth AuthConfig
}

type DbConfig struct {
	Dsn string
}

type AuthConfig struct {
	Secret string
}

func LoadConfig() *Config {
	err := godotenv.Load("..\\.env")
	if err != nil {
		log.Printf("Error loading .env file: %s. Using default config.", err.Error())
	}
	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("TOKEN"),
		},
	}
}
