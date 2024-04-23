package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (wrapper *apiWrapper) login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := loginRequestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	user := findUserByEmail(reqBody.Email, wrapper.users)
	if user == nil {
		msg := fmt.Sprintf("Could not find user with email: %s", reqBody.Email)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(reqBody.Password))
	if err != nil {
		msg := "Incorrect password"
		respondWithError(w, http.StatusUnauthorized, msg)
		return
	}

	currentTimeUTC := time.Now().UTC()
	tokenDuration := 900
	if reqBody.ExpiresInSeconds > 0 {
		tokenDuration = reqBody.ExpiresInSeconds
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(currentTimeUTC),
		ExpiresAt: jwt.NewNumericDate(currentTimeUTC.Add(time.Duration(tokenDuration) * time.Second)),
		Subject:   fmt.Sprint(user.Id),
	})
	signedToken, err := token.SignedString(wrapper.jwtSecret)
	if err != nil {
		msg := fmt.Sprintf("error signing token: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	userResponse := loginResponseBody{
		User: userResponse{
			Id:    user.Id,
			Email: user.Email,
		},
		Token: signedToken,
	}

	dat, err := json.Marshal(userResponse)
	if err != nil {
		msg := fmt.Sprintf("Error marshalling JSON response: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(dat)
}
