package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AIClient struct {
	APIKey string
}

func NewAIClient(apiKey string) *AIClient {
	return &AIClient{APIKey: apiKey}
}

func (c *AIClient) Ask(prompt string) (string, error) {
	const url = "https://openrouter.ai/api/v1/chat/completions"

	body := map[string]interface{}{
		"model": "mistralai/mistral-7b-instruct:free",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	jsonBody, _ := json.Marshal(body)

	// timeout 10 detik
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("HTTP-Referer", "https://agungperdana.store")
	req.Header.Set("X-Title", "AI Weather Assistant")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("AI timeout: request took too long")
		}
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// cek error dari API
	if errData, ok := result["error"]; ok {
		return "", fmt.Errorf("AI error: %v", errData)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid message")
	}

	text, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid content")
	}

	return text, nil
}
