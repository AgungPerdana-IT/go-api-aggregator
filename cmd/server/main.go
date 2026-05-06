package main

import (
	"log"
	"net/http"

	"github.com/AgungPerdana-IT/go-api-aggregator/internal/client"
	"github.com/AgungPerdana-IT/go-api-aggregator/internal/config"
	"github.com/AgungPerdana-IT/go-api-aggregator/internal/handler"
	"github.com/AgungPerdana-IT/go-api-aggregator/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	// init dependencies
	weatherClient := client.NewWeatherClient(cfg.WeatherAPIKey)
	weatherService := service.NewWeatherService(weatherClient)
	weatherHandler := handler.NewWeatherHandler(weatherService)

	aiClient := client.NewAIClient(cfg.AIAPIKey)
	aiService := service.NewAIService(aiClient)
	aiHandler := handler.NewAIHandler(aiService)

	http.HandleFunc("/api/ai", aiHandler.Ask)

	// routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/api/summary", weatherHandler.GetSummary(aiService))

	http.HandleFunc("/api/weather", weatherHandler.GetWeather)

	port := ":" + cfg.Port
	log.Println("Server running on", port)
	log.Println("Weather API Key:", cfg.WeatherAPIKey)
	log.Println("AI API Key:", cfg.AIAPIKey)

	log.Fatal(http.ListenAndServe(port, nil))
}
