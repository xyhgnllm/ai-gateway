package handler

import (
	"ai-gateway/internal/database"
	"ai-gateway/internal/response"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type CreateModelRequest struct {
	Name                  string `json:"name"`
	Provider              string `json:"provider"`
	InputPricePer1kCents  int64  `json:"input_price_per_1k_cents"`
	OutputPricePer1kCents int64  `json:"output_price_per_1k_cents"`
}

func (h *UserHandler) CreateModel(w http.ResponseWriter, r *http.Request) {
	var req CreateModelRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Provider = strings.TrimSpace(req.Provider)

	if req.Name == "" {
		response.Fail(w, http.StatusBadRequest, 400, "model name is required")
		return
	}

	if req.Provider == "" {
		req.Provider = "openai"
	}
	model, err := h.queries.CreateModel(r.Context(), database.CreateModelParams{
		Name:                  req.Name,
		Provider:              req.Provider,
		InputPricePer1kCents:  req.InputPricePer1kCents,
		OutputPricePer1kCents: req.OutputPricePer1kCents,
	})
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "create model failed")
		return
	}

	response.OK(w, model)

}

func (h *UserHandler) ListModles(w http.ResponseWriter, r *http.Request) {
	models, err := h.queries.ListModels(r.Context())
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "list models failed")
		return

	}

	response.OK(w, models)
}

type UpdateModelStatusRequest struct {
	Status string `json:"status"`
}

func (h *UserHandler) UpdataModelStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	modelID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid model id")
		return
	}

	var req UpdateModelStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid request body")
		return
	}

	req.Status = strings.TrimSpace(req.Status)
	if req.Status != "active" && req.Status != "disabled" {
		response.Fail(w, http.StatusBadRequest, 400, "invalid status")
		return
	}

	model, err := h.queries.UpdateModelStatus(r.Context(), database.UpdateModelStatusParams{
		ID:     modelID,
		Status: req.Status,
	})
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "update model status failed")
		return
	}

	response.OK(w, model)
}
