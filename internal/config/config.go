package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	AIAPIKey      string
	WeatherAPIKey string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		Port:          getEnv("PORT", "3000"),
		AIAPIKey:      os.Getenv("AI_API_KEY"),
		WeatherAPIKey: os.Getenv("WEATHER_API_KEY"),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
