package main

import (
	"fmt"
	"io"
	"net/http"
)

// Server struct holds the HTTP server and configuration
type Server struct {
	server *http.Server
}

// NewServer creates a new instance of Server with default settings
func NewServer() *Server {
	mux := http.NewServeMux()

	mux.Handle("/app/*", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.HandleFunc("/healthz", healthHandler)

	corsMux := addCORSHeaders(mux)

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

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

// addCORSHeaders is a custom middleware that adds CORS headers to the response
func addCORSHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
