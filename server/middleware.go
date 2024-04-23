package server

import "net/http"

func (wrapper *apiWrapper) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapper.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

// MiddleWareCORS is a custom middleware that adds CORS headers to the response
func MiddleWareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
