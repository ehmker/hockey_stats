package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_URL string
}

func Read() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db_url := os.Getenv("DB_URL")
	return Config{
		DB_URL: db_url,
	}
	
}