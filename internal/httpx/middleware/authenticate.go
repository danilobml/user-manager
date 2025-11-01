package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/danilobml/user-manager/internal/helpers"
	"github.com/danilobml/user-manager/internal/user/jwt"
)

type ctxKey string

const claimsCtxKey ctxKey = "claims"

func Authenticate(jwtManager *jwt.JwtManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			parts := strings.Fields(authHeader)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				helpers.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			tokenString := parts[1]

			claims, err := jwtManager.ParseAndValidateToken(tokenString)
			if err != nil {
				helpers.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			ctx := context.WithValue(r.Context(), claimsCtxKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helper
func GetClaimsFromContext(ctx context.Context) (*jwt.Claims, bool) {
	claims, ok := ctx.Value(claimsCtxKey).(*jwt.Claims)
	return claims, ok
}
