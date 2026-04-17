package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	SERVER_KEY     string
	AllowedOrigins []string
}

var AppConfig *Config

// LoadConfig baca file .env
func LoadConfig() {
	// load .env hanya sekali
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	AppConfig = &Config{
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		AllowedOrigins: []string{os.Getenv("DEV_ORIGIN"), os.Getenv("PROD_ORIGIN")},
	}
}
