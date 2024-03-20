package server

import (
	"fmt"
	"net/http"
)

// Server struct holds the HTTP server and configuration
type Server struct {
	server *http.Server
}

// NewServer creates a new instance of Server with default settings
func NewServer() *Server {
	wrapper := apiWrapper{fileServerHits: 0}

	mux := http.NewServeMux()
	mux.Handle("/app/*", wrapper.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", wrapper.healthHandler)
	mux.HandleFunc("GET /admin/metrics", wrapper.metricsHandler)
	mux.HandleFunc("POST /api/reset", wrapper.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", wrapper.validateHandler)

	corsMux := MiddleWareCORS(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	return &Server{server: server}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	fmt.Println("Server is running on http://localhost:8080")
	return s.server.ListenAndServe()
}
