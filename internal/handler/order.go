package handler

import (
	"ai-gateway/internal/database"
	"ai-gateway/internal/middleware"
	"ai-gateway/internal/response"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateOrderRequest struct {
	AmountCents int64 `json:"amount_cents"`
}

func generateOrderNo() (string, error) {
	bytes := make([]byte, 8)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "ORD" + time.Now().Format("20060102150405") + hex.EncodeToString(bytes), nil
}

func (h *UserHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)

	if !ok {
		response.Fail(w, http.StatusUnauthorized, 401, "unauthorized")
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid request body")
		return
	}

	if req.AmountCents <= 0 {
		response.Fail(w, http.StatusBadRequest, 400, "amount must be greater than 0")
		return
	}

	orderNo, err := generateOrderNo()
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "generate order No. failed")
		return
	}

	order, err := h.queries.CreateOrder(r.Context(), database.CreateOrderParams{
		UserID:      userID,
		OrderNo:     orderNo,
		AmountCents: req.AmountCents,
	})

	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "create order failed")
		return
	}

	response.OK(w, order)

}

func (h *UserHandler) ListMyOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, 401, "unauthorized")
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

	orders, err := h.queries.ListOrdersByUser(r.Context(), database.ListOrdersByUserParams{
		UserID: userID,
		Limit:  int32(page),
		Offset: int32(offset),
	})
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "list orders failed")
		return
	}

	response.OK(w, orders)

}
func (h *UserHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	pageSize, _ := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64)

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	orders, err := h.queries.ListOrders(r.Context(), database.ListOrdersParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "list orders failed")
		return
	}

	response.OK(w, orders)

}

func (h *UserHandler) PayOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, 400, "invalid order id")
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "begin transaction failed")
		return
	}
	defer tx.Rollback(r.Context())

	qtx := h.queries.WithTx(tx)

	order, err := qtx.GetOrderByID(r.Context(), orderID)
	if err != nil {
		response.Fail(w, http.StatusNotFound, 404, "order not found")
		return
	}

	if order.Status == "paid" {
		response.Fail(w, http.StatusBadRequest, 400, "order already paid")
		return
	}

	paidOrder, err := qtx.MarkOrderPaid(r.Context(), orderID)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "mark order paid failed")
		return
	}

	beforeUser, err := qtx.GetUserByID(r.Context(), paidOrder.UserID)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "get user balance failed")
		return
	}

	user, err := qtx.AddUserBalance(r.Context(), database.AddUserBalanceParams{
		ID:           paidOrder.UserID,
		BalanceCents: paidOrder.AmountCents,
	})
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "add user balance failed")
		return
	}

	transaction, err := qtx.CreateBalanceTransaction(r.Context(), database.CreateBalanceTransactionParams{
		UserID:             paidOrder.UserID,
		Type:               "recharge",
		AmountCents:        paidOrder.AmountCents,
		BalanceBeforeCents: beforeUser.BalanceCents,
		BalanceAfterCents:  user.BalanceCents,
		OrderID:            pgtype.Int8{Int64: paidOrder.ID, Valid: true},
		UsageLogID:         pgtype.Int8{Valid: false},
		Remark:             "order recharge",
	})
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "create balance transaction failed")
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "commit transaction failed")
		return
	}

	response.OK(w, map[string]any{
		"order":       paidOrder,
		"user":        user,
		"transaction": transaction,
	})
}
