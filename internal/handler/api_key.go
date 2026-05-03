package handler

import (
	"ai-gateway/internal/auth"
	"ai-gateway/internal/database"
	"ai-gateway/internal/middleware"
	"ai-gateway/internal/response"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

type APIKeyResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	KeyPrefix  string `json:"key_prefix"`
	Status     string `json:"status"`
	LastUsedAt any    `json:"last_used_at"`
	CreatedAt  any    `json:"created_at"`
}

func (h *UserHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, 401, "unauthorized")
		return
	}

	var req CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		req.Name = "default"
	}

	key, err := auth.GenerateAPIKey()
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "generate api key failed")
		return
	}
	keyHash := auth.HashAPIKey(key)
	keyPrefix := auth.APIKeyPrefix(key)

	apiKey, err := h.queries.CreateAPIKey(r.Context(), database.CreateAPIKeyParams{
		UserID:    userID,
		Name:      req.Name,
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
	})

	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "create api key failed")
		return
	}

	response.OK(w, map[string]any{
		"api_key": map[string]any{
			"key_prefix": apiKey.KeyPrefix,
		},
		"record": apiKey,
	})

}

func (h *UserHandler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, 401, "unauthorized")
		return
	}

	apiKeys, err := h.queries.ListAPIKeysByUser(r.Context(), userID)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "list api keys failed")
		return
	}

	result := make([]APIKeyResponse, 0, len(apiKeys))

	for _, key := range apiKeys {
		result = append(result, APIKeyResponse{
			ID:         key.ID,
			Name:       key.Name,
			KeyPrefix:  key.KeyPrefix,
			Status:     key.Status,
			LastUsedAt: key.LastUsedAt,
			CreatedAt:  key.CreatedAt,
		})
	}

	response.OK(w, result)

}

func (h *UserHandler) DisableAPIKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, 401, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")

	apiKeyID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid api key id")
		return
	}

	apiKey, err := h.queries.UpdateAPIKeyStatus(r.Context(), database.UpdateAPIKeyStatusParams{
		ID:     apiKeyID,
		Status: "disabled",
		UserID: userID,
	})

	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "disable api key failed")
		return
	}

	response.OK(w, apiKey)

}
