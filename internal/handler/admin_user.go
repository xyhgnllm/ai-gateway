package handler

import (
	"ai-gateway/internal/database"
	"ai-gateway/internal/response"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
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

	result := make([]UserResponse, 0, len(users))

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

type UpdateUserStatusRequest struct {
	Status string `json:"status"`
}

func (h *UserHandler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid user id")
		return
	}

	var req UpdateUserStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid request body")
		return
	}

	req.Status = strings.TrimSpace(req.Status)
	if req.Status != "active" && req.Status != "disabled" {
		response.Fail(w, http.StatusBadRequest, 400, "invalid status")
		return
	}

	user, err := h.queries.UpdateUserStatus(r.Context(), database.UpdateUserStatusParams{
		ID:     userID,
		Status: req.Status,
	})
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "update user status failed")
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

type UpdateUserRoleRequest struct {
	Role string `json:"role"`
}

func (h *UserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid user id")
		return
	}

	var req UpdateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid request body")
		return
	}

	req.Role = strings.TrimSpace(req.Role)
	if req.Role != "user" && req.Role != "admin" {
		response.Fail(w, http.StatusBadRequest, 400, "invalid role")
		return
	}

	user, err := h.queries.UpdateUserRole(r.Context(), database.UpdateUserRoleParams{
		ID:   userID,
		Role: req.Role,
	})
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "update user role failed")
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
