package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()

		if r.Header.Get("X-Request-ID") == "" {
			r.Header.Add("X-Request-ID", id)
		}

		next.ServeHTTP(w, r)
	})
}
