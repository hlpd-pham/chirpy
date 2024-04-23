package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hlpd-pham/chirpy/server/auth"
	"golang.org/x/crypto/bcrypt"
)

func (wrapper *apiWrapper) login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := userRequest{}
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

	accessToken, err := auth.GetSignedToken(auth.CHIRPY_ACCESS_ISSUER, fmt.Sprint(user.Id), wrapper.jwtSecret)
	if err != nil {
		msg := fmt.Sprintf("error signing access token: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	refreshToken, err := auth.GetSignedToken(auth.CHIRPY_REFRESH_ISSUER, fmt.Sprint(user.Id), wrapper.jwtSecret)
	if err != nil {
		msg := fmt.Sprintf("error signing refresh token: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	userResponse := loginResponseBody{
		Id:           user.Id,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
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
