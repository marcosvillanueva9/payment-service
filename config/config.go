package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL   string
	JWTSecret string
	PORT     string
	APP_ENV string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	return &Config{
		DBURL:    os.Getenv("DB_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		PORT:     os.Getenv("PORT"),
		APP_ENV:  os.Getenv("APP_ENV"),
	}, nil
}