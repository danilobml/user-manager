package httpx

import (
	"log"
	"net/http"
)

func Serve(port string, handler *http.Handler) {
	srv := http.Server{
		Addr: port,
		Handler: *handler,
	}

	log.Printf("server listening on port%s...", port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Server initialization failed", err)
	}
}
