package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	Port           string
	AccessSecret   string
	RefreshSecret  string
	RedisAddr      string
	FrontendOrigin string
	Domain         string
	Secure         string
}

func Load() (*Config, error) {
	var err error = godotenv.Load()

	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	var config = &Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		Port:           os.Getenv("PORT"),
		AccessSecret:   os.Getenv("ACCESS_SECRET"),
		RefreshSecret:  os.Getenv("REFRESH_SECRET"),
		RedisAddr:      os.Getenv("REDIS_ADDR"),
		FrontendOrigin: os.Getenv("FRONTEND_ORIGIN"),
		Domain:         os.Getenv("DOMAIN"),
		Secure:         os.Getenv("SECURE"),
	}

	return config, nil
}
