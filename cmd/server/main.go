package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	WeatherAPIKey string
	AIAPIKey      string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		Port:          os.Getenv("PORT"),
		WeatherAPIKey: os.Getenv("WEATHER_API_KEY"),
		AIAPIKey:      os.Getenv("AI_API_KEY"),
	}
}
