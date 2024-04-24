package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hlpd-pham/chirpy/server/auth"
	"golang.org/x/crypto/bcrypt"
)

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
	msg := fmt.Sprintf("error marshalling JSON response: %s", err)
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
	reqBody := userRequest{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		msg := fmt.Sprintf("error decoding request body: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	if findUserByEmail(reqBody.Email, wrapper.users) != nil {
		msg := "email is used"
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	email, err := validateUserEmail(reqBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), len(reqBody.Password))
	if err != nil {
		msg := fmt.Sprintf("error hashing password: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	newUser := user{
		Id:           wrapper.nextUserId,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	res := userResponse{
		Id:    newUser.Id,
		Email: newUser.Email,
	}

	dat, err := json.Marshal(res)
	if err != nil {
		msg := fmt.Sprintf("error marshalling JSON response: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	wrapper.users[wrapper.nextUserId] = newUser
	wrapper.nextUserId++

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write(dat)
}

func (wrapper *apiWrapper) updateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := userRequest{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		msg := fmt.Sprintf("error decoding request body: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	_, err = auth.GetToken(r, wrapper.jwtSecret, auth.CHIRPY_ACCESS_ISSUER)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user := findUserByEmail(reqBody.Email, wrapper.users)
	if user == nil {
		msg := fmt.Sprintf("could not find user with email: %s", reqBody.Email)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	if len(reqBody.Password) != 0 {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), len(reqBody.Password))
		if err != nil {
			msg := fmt.Sprintf("Error hashing password: %s", err)
			respondWithError(w, http.StatusServiceUnavailable, msg)
			return
		}
		user.PasswordHash = string(passwordHash)
	}

	res := userResponse{
		Id:    user.Id,
		Email: user.Email,
	}

	dat, err := json.Marshal(res)
	if err != nil {
		msg := fmt.Sprintf("Error marshalling JSON response: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(dat)
}
