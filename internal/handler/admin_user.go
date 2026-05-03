package handler

import (
	"ai-gateway/internal/database"
	"ai-gateway/internal/response"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type AddBalanceRequest struct {
	AmountCents int64 `json:"amount_cents"`
}

func (h *UserHandler) AddBalance(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "id")

	userID, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid user id")
		return
	}

	var req AddBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid request body")
		return
	}

	if req.AmountCents <= 0 {
		response.Fail(w, http.StatusBadRequest, 400, "amount must be greater than 0")
		return
	}

	user, err := h.queries.AddUserBalance(r.Context(), database.AddUserBalanceParams{
		ID:           userID,
		BalanceCents: req.AmountCents,
	})

	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "add balance failed")
		return
	}

	response.OK(w, UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Role:         user.Role,
		Status:       user.Status,
		BalanceCents: user.BalanceCents,
	})

}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	pageSize, _ := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64)

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	users, err := h.queries.ListUsers(r.Context(), database.ListUsersParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})

	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "list users failed")
		return
	}

	result := make([]UserResponse, 9, len(users))

	for _, user := range users {
		result = append(result, UserResponse{
			ID:           user.ID,
			Email:        user.Email,
			Name:         user.Name,
			Role:         user.Role,
			Status:       user.Status,
			BalanceCents: user.BalanceCents,
		})
	}

	response.OK(w, result)
}
