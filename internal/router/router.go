package router

import (
	"ai-gateway/internal/database"
	"ai-gateway/internal/handler"
	"net/http"

	appmiddleware "ai-gateway/internal/middleware"

	"github.com/go-chi/chi"
)

func New(
	queries *database.Queries,
	jwtSecret string,
	openAIBaseURL string,
	openAIAPIKey string,
) http.Handler {

	r := chi.NewRouter()

	r.Use(appmiddleware.Common)

	userHandler := handler.NewUserHandler(queries, jwtSecret)
	gatewayHandler := handler.NewGatewayHandler(queries, openAIBaseURL, openAIAPIKey)

	r.Get("/health", handler.Health)
	r.Post("/users", userHandler.Create)
	r.Post("/login", userHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(appmiddleware.JWT(jwtSecret))

		r.Get("/me", userHandler.Me)
		r.Post("/api-keys", userHandler.CreateAPIKey)
		r.Get("/api-keys", userHandler.ListAPIKeys)
		r.Get("/v1/test", gatewayHandler.GatewayTest)
		r.Post("/v1/chat/completions", gatewayHandler.ChatCompletions)
		r.Get("/usage-logs", userHandler.ListUsageLogs)
		r.Patch("/api-keys/{id}/disable", userHandler.DisableAPIKey)
		r.Post("/orders", userHandler.CreateOrder)
		r.Get("/orders", userHandler.ListMyOrders)
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
	})

	return r
}
