package service

import "github.com/AgungPerdana-IT/go-api-aggregator/internal/client"

type AIService struct {
	Client *client.AIClient
}

func NewAIService(c *client.AIClient) *AIService {
	return &AIService{Client: c}
}

func (s *AIService) Ask(prompt string) (string, error) {
	return s.Client.Ask(prompt)
}
