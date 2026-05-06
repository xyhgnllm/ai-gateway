package router

import (
	"ai-gateway/internal/database"
	"ai-gateway/internal/handler"
	"net/http"
	"time"

	appmiddleware "ai-gateway/internal/middleware"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(
	queries *database.Queries,
	db *pgxpool.Pool,
	jwtSecret string,
	openAIBaseURL string,
	openAIAPIKey string,
) http.Handler {

	r := chi.NewRouter()

	r.Use(appmiddleware.Common)

	userHandler := handler.NewUserHandler(queries, db, jwtSecret)
	gatewayHandler := handler.NewGatewayHandler(queries, db, openAIBaseURL, openAIAPIKey)

	r.Get("/health", handler.Health)
	r.Post("/users", userHandler.Create)
	r.Post("/login", userHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(appmiddleware.JWT(jwtSecret))
		r.Use(appmiddleware.RateLimit(60, time.Minute))

		r.Get("/me", userHandler.Me)
		r.Post("/api-keys", userHandler.CreateAPIKey)
		r.Get("/api-keys", userHandler.ListAPIKeys)
		r.Get("/v1/test", gatewayHandler.GatewayTest)
		r.Post("/v1/chat/completions", gatewayHandler.ChatCompletions)
		r.Get("/usage-logs", userHandler.ListUsageLogs)
		r.Patch("/api-keys/{id}/disable", userHandler.DisableAPIKey)
		r.Post("/orders", userHandler.CreateOrder)
		r.Get("/orders", userHandler.ListMyOrders)
		r.Get("/balance-transactions", userHandler.ListBalanceTransactions)
	})

	r.Group(func(r chi.Router) {
		r.Use(appmiddleware.JWT(jwtSecret))
		r.Use(appmiddleware.RequireAdmin)

		r.Post("/admin/users/{id}/balance", userHandler.AddBalance)
		r.Get("/admin/users", userHandler.ListUsers)
		r.Post("/admin/models", userHandler.CreateModel)
		r.Get("/admin/models", userHandler.ListModles)
		r.Patch("/admin/models/{id}/status", userHandler.UpdataModelStatus)
		r.Post("/admin/orders/{id}/pay", userHandler.PayOrder)
		r.Get("/admin/users/{id}/balance-transactions", userHandler.ListUserBalanceTransactions)
		r.Get("/admin/users/{id}/usage-logs", userHandler.ListUserUsageLogs)
		r.Patch("/admin/users/{id}/status", userHandler.UpdateUserStatus)
		r.Patch("/admin/users/{id}/role", userHandler.UpdateUserRole)
		r.Get("/admin/stats", userHandler.AdminStats)
		r.Get("/admin/stats/models", userHandler.ModelUsageStats)
		r.Get("/admin/stats/users", userHandler.UserUsageStats)
		r.Get("/admin/stats/daily", userHandler.DailyUsageStats)
	})

	return r
}
