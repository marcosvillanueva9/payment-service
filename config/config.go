package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPATH   string
	JWTSecret string
	PORT     string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	return &Config{
		DBPATH:   os.Getenv("DB_PATH"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		PORT:     os.Getenv("PORT"),
	}, nil
}