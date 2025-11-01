package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := statusWriter{ResponseWriter: w}

		next.ServeHTTP(&sw, r)

		if sw.status == 0 {
			sw.status = http.StatusOK
		}

		requestId := r.Header.Get("X-Request-ID") 

		log.Printf("Request URI: %s, RequestId: %s, Method: %s, Status: %d, Latency: %v", r.RequestURI, requestId, r.Method, sw.status,time.Since(start))
	})
}
