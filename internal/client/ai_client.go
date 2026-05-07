package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AIClient struct {
	APIKey string
}

func NewAIClient(apiKey string) *AIClient {
	return &AIClient{APIKey: apiKey}
}

func (c *AIClient) Ask(prompt string) (string, error) {
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-flash-latest:generateContent?key=%s",
		c.APIKey,
	)

	fmt.Println("Tembak URL:", url)

	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// 🔥 DEBUG (sementara)
	fmt.Println("AI RAW:", result)

	// cek error dari API
	if errData, ok := result["error"]; ok {
		return "", fmt.Errorf("AI error: %v", errData)
	}

	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("no candidates returned")
	}

	content, ok := candidates[0].(map[string]interface{})["content"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid content")
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return "", fmt.Errorf("invalid parts")
	}

	text, ok := parts[0].(map[string]interface{})["text"].(string)
	if !ok {
		return "", fmt.Errorf("invalid text")
	}

	return text, nil
}
