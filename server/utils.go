package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func validateUserEmail(reqBody userRequest) (string, error) {
	if len(reqBody.Email) == 0 {
		return "", errors.New("email is required to create account")
	}
	return reqBody.Email, nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	log.Println(msg)
	respBody := errorResponse{
		Message: msg,
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

func findUserByEmail(email string, users map[int]user) *user {
	for _, user := range users {
		if user.Email == email {
			return &user
		}
	}
	return nil
}
