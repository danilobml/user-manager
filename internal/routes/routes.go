package routes

import (
	"encoding/json"
	"net/http"

	"github.com/danilobml/user-manager/internal/httpx/middleware"
	"github.com/danilobml/user-manager/internal/user/handler"
)

func NewRouter(userHandler *handler.UserHandler, authMiddleware middleware.Middleware) http.Handler {
	mux := http.NewServeMux()

	// Public
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)
	mux.HandleFunc("POST /request-password", userHandler.RequestPasswordReset)

	// Protected
	mux.Handle("DELETE /users/{id}",
		authMiddleware(http.HandlerFunc(userHandler.UnregisterUser)),
	)
	mux.Handle("PUT /users/{id}",
		authMiddleware(http.HandlerFunc(userHandler.UpdateUser)),
	)
	mux.Handle("PUT /users/reset-password",
		authMiddleware(http.HandlerFunc(userHandler.ResetPassword)),
	)
	// Admin
	mux.Handle("GET /users",
		authMiddleware(http.HandlerFunc(userHandler.GetAllUsers)),
	)
	mux.Handle("DELETE /users/{id}/remove",
		authMiddleware(http.HandlerFunc(userHandler.RemoveUser)),
	)

	// Global middlewares
	use := middleware.ApplyMiddlewares(
		middleware.Recover,
		middleware.Cors,
		middleware.RequestId,
		middleware.Logger,
	)

	return use(mux)
}

func health(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"health": "ok"}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
