package handler

import (
	"ai-gateway/internal/auth"
	"ai-gateway/internal/database"
	"ai-gateway/internal/middleware"
	"ai-gateway/internal/response"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserHandler struct {
	queries   *database.Queries
	jwtSecret string
}

func NewUserHandler(queries *database.Queries, jwtSecret string) *UserHandler {
	return &UserHandler{
		queries:   queries,
		jwtSecret: jwtSecret,
	}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserResponse struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	Status       string `json:"status"`
	BalanceCents int64  `json:"balance_cents"`
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invaild request body")
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	req.Name = strings.TrimSpace(req.Name)

	if req.Email == "" {
		response.Fail(w, http.StatusBadRequest, 400, "email is required")
		return
	}

	if len(req.Password) < 8 {
		response.Fail(w, http.StatusBadRequest, 400, "password must be at least 8 characters")
		return
	}

	if req.Name == "" {
		response.Fail(w, http.StatusBadRequest, 400, "name is required")
		return
	}
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "hash password failed")
		return
	}

	user, err := h.queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Name:         req.Name,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			response.Fail(w, http.StatusBadRequest, 400, "email already exists")
			return
		}

		response.Fail(w, http.StatusInternalServerError, 500, "create user failed")
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid request body")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" {
		response.Fail(w, http.StatusBadRequest, 400, "email is required")
		return
	}

	if req.Password == "" {
		response.Fail(w, http.StatusBadRequest, 400, "password is required")
		return
	}

	user, err := h.queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		response.Fail(w, http.StatusUnauthorized, 401, "invalid email or password")
		return
	}
	if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
		response.Fail(w, http.StatusUnauthorized, 401, "invalid email or password")
		return
	}
	token, err := auth.GenerateToken(user.ID, user.Role, h.jwtSecret)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "generate token failed")
		return
	}

	response.OK(w, map[string]any{
		"user": UserResponse{
			ID:           user.ID,
			Email:        user.Email,
			Name:         user.Name,
			Role:         user.Role,
			Status:       user.Status,
			BalanceCents: user.BalanceCents},
		"token": token,
	})

}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, 401, "unauthorized")
		return
	}

	user, err := h.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "get user failed")
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
