package handler

import (
	"ai-gateway/internal/database"
	"ai-gateway/internal/middleware"
	"ai-gateway/internal/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *UserHandler) ListUsageLogs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)

	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	pageSize, _ := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64)

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	if !ok {
		response.Fail(w, http.StatusUnauthorized, 401, "unauthorized")
		return
	}

	logs, err := h.queries.ListUsageLogsByUser(r.Context(), database.ListUsageLogsByUserParams{
		UserID: userID,
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})

	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "list usage logs failed")
		return
	}

	response.OK(w, logs)
}

func (h *UserHandler) ListUserUsageLogs(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid user id")
		return
	}

	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	pageSize, _ := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64)

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	logs, err := h.queries.ListUsageLogsByUser(r.Context(), database.ListUsageLogsByUserParams{
		UserID: userID,
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "list user usage logs failed")
		return
	}

	response.OK(w, logs)
}
