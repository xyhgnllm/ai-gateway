package middleware

import (
	"ai-gateway/internal/auth"
	"ai-gateway/internal/database"
	"ai-gateway/internal/response"
	"context"
	"net/http"
	"strings"
)

func APIKey(queries *database.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Fail(w, http.StatusUnauthorized, 401, "missing api key")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Fail(w, http.StatusUnauthorized, 401, "invalid authorization header")
				return
			}

			apiKey := parts[1]
			KeyHash := auth.HashAPIKey(apiKey)

			record, err := queries.GetAPIKeyByHash(r.Context(), KeyHash)
			if err != nil {
				response.Fail(w, http.StatusUnauthorized, 401, "invalid api key")
				return
			}

			if record.Status != "active" {
				response.Fail(w, http.StatusUnauthorized, 401, "api key is disabled")
				return
			}

			_ = queries.UpdateAPIKeyLastUsed(r.Context(), record.ID)

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, record.UserID)
			ctx = context.WithValue(ctx, APIKeyIDKey, record.ID)

			next.ServeHTTP(w, r.WithContext(ctx))

		})

	}
}
