package handler

import (
	"ai-gateway/internal/auth"
	"ai-gateway/internal/database"
	"ai-gateway/internal/middleware"
	"ai-gateway/internal/response"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
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

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
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
		"api_key": key,
		"record":  apiKey,
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

	response.OK(w, apiKeys)

}

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
