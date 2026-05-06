package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AgungPerdana-IT/go-api-aggregator/internal/service"
)

type WeatherHandler struct {
	Service *service.WeatherService
}

func NewWeatherHandler(s *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{Service: s}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "city is required", http.StatusBadRequest)
		return
	}

	data, err := h.Service.GetWeather(city)
	if err != nil {
		http.Error(w, "failed to fetch weather", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *WeatherHandler) GetSummary(aiService *service.AIService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		city := r.URL.Query().Get("city")
		if city == "" {
			http.Error(w, "city is required", http.StatusBadRequest)
			return
		}

		data, err := h.Service.GetSummary(city, aiService)
		if err != nil {
			http.Error(w, "failed to get summary", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(data)
	}
}
