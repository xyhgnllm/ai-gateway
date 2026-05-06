package handler

import (
	"net/http"

	"ai-gateway/internal/response"
)

func (h *UserHandler) AdminStats(w http.ResponseWriter, r *http.Request) {
	userCount, err := h.queries.CountUsers(r.Context())
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "count users failed")
		return
	}

	orderCount, err := h.queries.CountOrders(r.Context())
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "count orders failed")
		return
	}

	paidAmount, err := h.queries.SumPaidOrderAmount(r.Context())
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "sum paid order amount failed")
		return
	}

	usageCount, err := h.queries.CountUsageLogs(r.Context())
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "count usage logs failed")
		return
	}

	usageCost, err := h.queries.SumUsageCost(r.Context())
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "sum usage cost failed")
		return
	}

	response.OK(w, map[string]any{
		"user_count":        userCount,
		"order_count":       orderCount,
		"paid_amount_cents": paidAmount,
		"usage_count":       usageCount,
		"usage_cost_cents":  usageCost,
	})
}

func (h *UserHandler) ModelUsageStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.queries.ListModelUsageStats(r.Context())
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "list model usage stats failed")
		return
	}

	response.OK(w, stats)
}

func (h *UserHandler) UserUsageStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.queries.ListUserUsageStats(r.Context())
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "list user usage stats failed")
		return
	}

	response.OK(w, stats)
}

func (h *UserHandler) DailyUsageStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.queries.ListDailyUsageStats(r.Context())
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, 500, "list daily usage stats failed")
		return
	}

	response.OK(w, stats)
}
