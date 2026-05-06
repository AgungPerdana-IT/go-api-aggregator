package service

import (
	"fmt"

	"github.com/AgungPerdana-IT/go-api-aggregator/internal/client"
)

type WeatherService struct {
	Client *client.WeatherClient
}

func NewWeatherService(c *client.WeatherClient) *WeatherService {
	return &WeatherService{Client: c}
}

func (s *WeatherService) GetWeather(city string) (*WeatherResponse, error) {
	data, err := s.Client.GetWeather(city)
	if err != nil {
		return nil, err
	}

	// cek kalau API balikin error
	if code, ok := data["cod"].(float64); ok && code != 200 {
		return nil, fmt.Errorf("weather API error: %v", data["message"])
	}

	main, ok := data["main"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid main data")
	}

	weatherArr, ok := data["weather"].([]interface{})
	if !ok || len(weatherArr) == 0 {
		return nil, fmt.Errorf("invalid weather data")
	}

	weather, ok := weatherArr[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid weather object")
	}

	return &WeatherResponse{
		City:        data["name"].(string),
		Temperature: main["temp"].(float64),
		Description: weather["description"].(string),
	}, nil
}

type WeatherResponse struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
}

func (s *WeatherService) GetSummary(city string, aiService *AIService) (map[string]interface{}, error) {
	weather, err := s.GetWeather(city)
	if err != nil {
		return nil, err
	}

	// bikin prompt ke AI
	prompt := "Cuaca di " + weather.City +
		" saat ini " + weather.Description +
		" dengan suhu " + fmt.Sprintf("%.1f", weather.Temperature) +
		"°C. Berikan saran singkat apakah cocok keluar rumah."

	advice, err := aiService.Ask(prompt)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"city":        weather.City,
		"temperature": weather.Temperature,
		"description": weather.Description,
		"advice":      advice,
	}, nil
}
