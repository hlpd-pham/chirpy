package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

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

func (wrapper *apiWrapper) getAllChirps(w http.ResponseWriter, _ *http.Request) {
	respBody := getAllChirpsResponse{
		Body: wrapper.chirps,
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

func (wrapper *apiWrapper) getOneChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if len(id) == 0 {
		respondWithError(w, http.StatusNotFound, "chirp ID is required to get information")
		return
	}

	chirpId, err := strconv.Atoi(id)
	if err != nil || chirpId < 0 || chirpId > len(wrapper.chirps) {
		msg := fmt.Sprintf("chirp ID: %s is invalid, err: %v", id, err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	dat, err := json.Marshal(wrapper.chirps[chirpId-1])
	msg := fmt.Sprintf("Error marshalling JSON response: %s", err)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(dat)
}

func (wrapper *apiWrapper) createChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := createChirpRequestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	words, err := validateChirp(reqBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respBody := createChirpResponse{
		Body: strings.Join(words, " "),
		Id:   wrapper.nextChirpId,
	}

	dat, err := json.Marshal(respBody)
	msg := fmt.Sprintf("Error marshalling JSON response: %s", err)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	wrapper.chirps = append(wrapper.chirps, chirp{Id: wrapper.nextChirpId, Body: respBody.Body})
	wrapper.nextChirpId++

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write(dat)
}

func validateChirp(reqBody createChirpRequestBody) ([]string, error) {
	if len(reqBody.Body) > 140 || len(reqBody.Body) == 0 {
		msg := fmt.Sprintf("chirp too or empty long, %d in length", len(reqBody.Body))
		return nil, errors.New(msg)
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
	return words, nil
}

func (wrapper *apiWrapper) getOneUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if len(id) == 0 {
		respondWithError(w, http.StatusNotFound, "user ID is required to get information")
		return
	}

	userId, err := strconv.Atoi(id)
	if err != nil || userId < 0 || userId > len(wrapper.chirps) {
		msg := fmt.Sprintf("user ID: %s is invalid, err: %v", id, err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	dat, err := json.Marshal(wrapper.users[userId-1])
	msg := fmt.Sprintf("Error marshalling JSON response: %s", err)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(dat)
}

func (wrapper *apiWrapper) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := createUserRequestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	email, err := validateUserEmail(reqBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	newUser := user{
		Id:    wrapper.nextUserId,
		Email: email,
	}

	dat, err := json.Marshal(newUser)
	msg := fmt.Sprintf("Error marshalling JSON response: %s", err)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	wrapper.nextUserId++
	wrapper.users = append(wrapper.users, newUser)

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write(dat)
}

func validateUserEmail(reqBody createUserRequestBody) (string, error) {
	if len(reqBody.Email) == 0 {
		return "", errors.New("email is required to create account")
	}
	return reqBody.Email, nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	log.Println(msg)
	respBody := errorResponse{
		Error: errors.New(msg),
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
