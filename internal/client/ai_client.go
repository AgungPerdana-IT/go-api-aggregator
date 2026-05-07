if c.APIKey == "" {
	return "", fmt.Errorf("API key is empty")
}

resp, err := http.DefaultClient.Do(req)
if err != nil {
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("AI timeout")
	}
	return "", err
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
	return "", fmt.Errorf("API error: status %d", resp.StatusCode)
}

var result map[string]interface{}
if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	return "", err
}

// handle error dari OpenRouter
if errData, ok := result["error"].(map[string]interface{}); ok {
	return "", fmt.Errorf("AI error: %v", errData["message"])
}

choices, ok := result["choices"].([]interface{})
if !ok || len(choices) == 0 {
	return "", fmt.Errorf("no choices returned: %v", result)
}

msg, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
if !ok {
	return "", fmt.Errorf("invalid message format")
}

text, ok := msg["content"].(string)
if !ok {
	return "", fmt.Errorf("invalid content")
}