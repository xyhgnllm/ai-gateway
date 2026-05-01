package router

import (
	"ai-gateway/internal/database"
	"ai-gateway/internal/handler"
	"net/http"

	appmiddleware "ai-gateway/internal/middleware"

	"github.com/go-chi/chi"
)

func New(queries *database.Queries, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	r.Use(appmiddleware.Common)

	userHandler := handler.NewUserHandler(queries, jwtSecret)

	r.Get("/health", handler.Health)
	r.Post("/users", userHandler.Create)
	r.Post("/login", userHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(appmiddleware.JWT(jwtSecret))

		r.Get("/me", userHandler.Me)
	})

	return r
}
