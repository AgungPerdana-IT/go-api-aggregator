package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WeatherClient struct {
	APIKey string
}

func NewWeatherClient(apiKey string) *WeatherClient {
	return &WeatherClient{APIKey: apiKey}
}

func (c *WeatherClient) GetWeather(city string) (map[string]interface{}, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric",
		city,
		c.APIKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
