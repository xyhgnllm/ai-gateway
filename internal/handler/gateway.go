package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"ai-gateway/internal/database"
	"ai-gateway/internal/middleware"
	"ai-gateway/internal/response"

	"github.com/jackc/pgx/v5/pgtype"
)

type GatewayHandler struct {
	queries       *database.Queries
	openAIBaseURL string
	openAIAPIKey  string
}

func NewGatewayHandler(queries *database.Queries, openAIBaseURL string, openAIAPIKey string) *GatewayHandler {
	return &GatewayHandler{
		queries:       queries,
		openAIBaseURL: openAIBaseURL,
		openAIAPIKey:  openAIAPIKey,
	}
}

func (h *GatewayHandler) GatewayTest(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, 401, "unauthorized")
		return
	}

	response.OK(w, map[string]any{
		"message": "api key valid",
		"user_id": userID,
	})
}

type ChatCompletionRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type ChatCompletionPayload struct {
	Model string `json:"model"`
}

func (h *GatewayHandler) ChatCompletions(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, 401, "unauthorized")
		return
	}

	apiKeyID, ok := r.Context().Value(middleware.APIKeyIDKey).(int64)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, 401, "missing api key id")
		return
	}

	user, err := h.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		response.Fail(w, http.StatusUnauthorized, 401, "user not found")
		return
	}

	if user.BalanceCents <= 0 {
		response.Fail(w, http.StatusPaymentRequired, 402, "insufficient balance")
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "read request body failed")
		return
	}

	var payload ChatCompletionPayload
	_ = json.Unmarshal(requestBody, &payload)

	payload.Model = strings.TrimSpace(payload.Model)

	if payload.Model == "" {
		response.Fail(w, http.StatusBadRequest, 400, "model is required")
		return
	}
	model, err := h.queries.GetActiveModelByName(r.Context(), payload.Model)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "model not available")
		return
	}

	_ = model

	upstreamURL := h.openAIBaseURL + "/v1/chat/completions"

	upstreamReq, err := http.NewRequestWithContext(
		r.Context(),
		http.MethodPost,
		upstreamURL,
		bytes.NewReader(requestBody),
	)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "create upstream request failed")
		return
	}

	upstreamReq.Header.Set("Content-Type", "application/json")
	upstreamReq.Header.Set("Authorization", "Bearer "+h.openAIAPIKey)

	upstreamResp, err := http.DefaultClient.Do(upstreamReq)
	if err != nil {
		response.Fail(w, http.StatusBadGateway, 502, "upstream request failed")
		return
	}
	defer upstreamResp.Body.Close()

	responseBody, err := io.ReadAll(upstreamResp.Body)
	if err != nil {
		response.Fail(w, http.StatusBadGateway, 502, "read upstream response failed")
		return
	}

	costCents := int64(1)

	if upstreamResp.StatusCode >= 200 && upstreamResp.StatusCode < 300 {
		_, err = h.queries.DeductUserBalance(r.Context(), database.DeductUserBalanceParams{
			ID:           userID,
			BalanceCents: costCents,
		})

		if err != nil {
			response.Fail(w, http.StatusPaymentRequired, 402, "insufficient balance")
			return
		}
	}

	_, _ = h.queries.CreateUsageLog(r.Context(), database.CreateUsageLogParams{
		UserID:            userID,
		ApiKeyID:          pgtype.Int8{Int64: apiKeyID, Valid: true},
		Model:             payload.Model,
		Endpoint:          "/v1/chat/completions",
		StatusCode:        int32(upstreamResp.StatusCode),
		Success:           upstreamResp.StatusCode >= 200 && upstreamResp.StatusCode < 300,
		RequestBodyBytes:  int64(len(requestBody)),
		ResponseBodyBytes: int64(len(responseBody)),
		CostCents:         costCents,
	})

	w.Header().Set("Content-Type", upstreamResp.Header.Get("Content-Type"))
	w.WriteHeader(upstreamResp.StatusCode)
	w.Write(responseBody)
}
