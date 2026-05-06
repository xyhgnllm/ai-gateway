package middleware

import (
	"net/http"
	"sync"
	"time"

	"ai-gateway/internal/response"
)

type visitor struct {
	count     int
	resetTime time.Time
}

var (
	visitors = make(map[int64]*visitor)
	mu       sync.Mutex
)

func RateLimit(maxRequests int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKeyID, ok := r.Context().Value(APIKeyIDKey).(int64)
			if !ok {
				response.Fail(w, http.StatusUnauthorized, 401, "missing api key id")
				return
			}

			now := time.Now()

			mu.Lock()
			v, exists := visitors[apiKeyID]
			if !exists || now.After(v.resetTime) {
				visitors[apiKeyID] = &visitor{
					count:     1,
					resetTime: now.Add(window),
				}
				mu.Unlock()

				next.ServeHTTP(w, r)
				return
			}

			if v.count >= maxRequests {
				mu.Unlock()
				response.Fail(w, http.StatusTooManyRequests, 429, "too many requests")
				return
			}

			v.count++
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}
