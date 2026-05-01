package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort   string
	DBURL     string
	JWTSercet string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found,using system env")
	}

	return &Config{
		AppPort:   getEnv("APP_PORT", "8080"),
		DBURL:     getEnv("DB_URL", ""),
		JWTSercet: getEnv("JWT_SECRET", "dev-secret"),
	}
}

func getEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return value
}
