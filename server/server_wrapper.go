package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type metricsData struct {
	HitCount int
}

type apiWrapper struct {
	fileServerHits int
}

type requestBody struct {
	Body string `json:"body"`
}

type response struct {
	Error       error  `json:"error"`
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

func (wrapper *apiWrapper) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
	`, wrapper.fileServerHits)
}

func (wrapper *apiWrapper) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapper.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

func (wrapper *apiWrapper) resetHandler(w http.ResponseWriter, r *http.Request) {
	wrapper.fileServerHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (wrapper *apiWrapper) healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

func (wrapper *apiWrapper) validateHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	if len(reqBody.Body) > 140 || len(reqBody.Body) == 0 {
		msg := fmt.Sprintf("chirp too or empty long, %d in length", len(reqBody.Body))
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	log.Printf("valid message has %d characters", len(reqBody.Body))
	words := strings.Split(reqBody.Body, " ")
	for idx, word := range words {
		for _, badWord := range []string{"kerfuffle", "sharbert", "fornax"} {
			if word == badWord {
				words[idx] = "****"
				break
			}
		}
	}

	respBody := response{
		Error:       nil,
		Valid:       true,
		CleanedBody: strings.Join(words, " "),
	}
	dat, err := json.Marshal(respBody)
	msg := fmt.Sprintf("Error marshalling JSON response: %s", err)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	log.Println(msg)
	respBody := response{
		Error: errors.New(msg),
		Valid: false,
	}
	dat, err := json.Marshal(respBody)
	msg = fmt.Sprintf("Error marshalling JSON response: %s", err)
	if err != nil {
		log.Println(msg)
		return
	}

	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	w.Write(dat)
}
