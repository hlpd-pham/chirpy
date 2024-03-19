package server

import (
	"io"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

func MetricsHandler(w http.ResponseWriter, _ *http.Request) {
}
