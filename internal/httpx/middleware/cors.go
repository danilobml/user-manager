package middleware

import (
	"slices"
	"net/http"

	"github.com/rs/cors"
)


func Cors(mux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := []string{"*"}

		origin := r.Header.Get("Origin")

		if origin != "" && !isInAllowedOrigins(allowedOrigins, origin) {
			http.Error(w, "CORS origin denied", http.StatusForbidden)
			return
		}
		
		c := cors.New(cors.Options{
			AllowedOrigins:   allowedOrigins,
			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Authorization"},
			AllowCredentials: false,
		})
		c.Handler(mux).ServeHTTP(w, r)
	})
}

// Helper:
func isInAllowedOrigins(allowedOrigins []string, checkOrigin string) bool {
	if slices.Contains(allowedOrigins, "*") {
		return true
	}
	return slices.Contains(allowedOrigins, checkOrigin)
}
