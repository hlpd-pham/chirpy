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
	wrapper := apiWrapper{
		fileServerHits: 0,
		nextChirpId:    1,
		nextUserId:     1,
		chirps:         []chirp{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/*", wrapper.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(".")))),
	)
	mux.HandleFunc("GET /api/healthz", wrapper.healthHandler)
	mux.HandleFunc("GET /admin/metrics", wrapper.metricsHandler)
	mux.HandleFunc("POST /api/reset", wrapper.resetHandler)

	mux.HandleFunc("POST /api/chirps", wrapper.createChirp)
	mux.HandleFunc("GET /api/chirps", wrapper.getAllChirps)
	mux.HandleFunc("GET /api/chirps/{id}", wrapper.getOneChirp)

	mux.HandleFunc("POST /api/users", wrapper.createUser)
	mux.HandleFunc("GET /api/users/{id}", wrapper.getOneUser)

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
