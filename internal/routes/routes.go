package routes

import (
	"encoding/json"
	"net/http"

	"github.com/danilobml/user-manager/internal/httpx/middleware"
	"github.com/danilobml/user-manager/internal/user/handler"
)

func NewRouter(userHandler *handler.UserHandler) http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)

	//TODO: this will be a protected route for admins
	mux.HandleFunc("GET /users", userHandler.GetAllUsers)

	useMiddlewares := middleware.ApplyMiddlewares(
		middleware.Recover,
		middleware.Cors,
		middleware.RequestId,
		middleware.Logger,
	)

	return useMiddlewares(mux)
}

func health(w http.ResponseWriter, r *http.Request) {
	resp :=  map[string]string{"health": "ok"}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
