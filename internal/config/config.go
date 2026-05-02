package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	DBURL         string
	JWTSecret     string
	OpenAIBaseURL string
	OpenAIAPIKey  string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found,using system env")
	}

	return &Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		DBURL:         getEnv("DB_URL", ""),
		JWTSecret:     getEnv("JWT_SECRET", "dev-secret"),
		OpenAIBaseURL: getEnv("OPENAI_BASE_URL", "https://api.openai.com"),
		OpenAIAPIKey:  getEnv("OPENAI_API_KEY", ""),
	}
}

func getEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return value
}
