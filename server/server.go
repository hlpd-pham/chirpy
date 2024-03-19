package server

import (
	"fmt"
	"net/http"
)

// Server struct holds the HTTP server and configuration
type Server struct {
	server *http.Server
}

type apiConfig struct {
	fileServerHits int
}

// NewServer creates a new instance of Server with default settings
func NewServer() *Server {
	apiCfg := apiConfig{fileServerHits: 0}

	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", HealthHandler)
	mux.HandleFunc("GET /api/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.handlerReset)

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

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", 2)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
