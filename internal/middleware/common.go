package middleware

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
)

func Common(next http.Handler) http.Handler {
	return middleware.Logger(
		middleware.Recoverer(next),
	)
}
