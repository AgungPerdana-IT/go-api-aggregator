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
	if c.APIKey == "" {
		return "", fmt.Errorf("API key is empty")
	}

	url := "https://openrouter.ai/api/v1/chat/completions"

	body := map[string]interface{}{
		"model": "mistralai/mistral-7b-instruct:free",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// timeout (naikin dikit biar ga gampang timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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

	// cek status code dulu (PENTING)
	if resp.StatusCode != http.StatusOK {
		var errBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errBody)
		return "", fmt.Errorf("API error: status %d, body: %v", resp.StatusCode, errBody)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// handle error dari OpenRouter
	if errData, ok := result["error"].(map[string]interface{}); ok {
		return "", fmt.Errorf("AI error: %v", errData["message"])
	}

	// parsing response
	choicesRaw, exists := result["choices"]
	if !exists {
		return "", fmt.Errorf("no choices field in response: %v", result)
	}

	choices, ok := choicesRaw.([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("invalid choices format")
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid choice format")
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid message format")
	}

	text, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid content format")
	}

	return text, nil
}
