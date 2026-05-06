package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AgungPerdana-IT/go-api-aggregator/internal/service"
)

type AIHandler struct {
	Service *service.AIService
}

func NewAIHandler(s *service.AIService) *AIHandler {
	return &AIHandler{Service: s}
}

type AIRequest struct {
	Prompt string `json:"prompt"`
}

func (h *AIHandler) Ask(w http.ResponseWriter, r *http.Request) {
	var req AIRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Prompt == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	res, err := h.Service.Ask(req.Prompt)
	if err != nil {
		http.Error(w, "failed to call AI", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"response": res,
	})
}
