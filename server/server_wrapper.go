package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
	reqBody := userRequest{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body: %s", err)
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
		msg := fmt.Sprintf("Error hashing password: %s", err)
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
		msg := fmt.Sprintf("Error marshalling JSON response: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	wrapper.nextUserId++
	wrapper.users = append(wrapper.users, newUser)

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write(dat)
}

func (wrapper *apiWrapper) updateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := userRequest{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		msg := "Could not find auth header from request"
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		msg := "Auth header is not formatted correctly"
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	user := findUserByEmail(reqBody.Email, wrapper.users)
	if user == nil {
		msg := fmt.Sprintf("Could not find user with email: %s", reqBody.Email)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	fmt.Printf("headers: %v\n", headerParts[1])

	token, err := jwt.ParseWithClaims(
		headerParts[1],
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return wrapper.jwtSecret, nil
		},
	)
	if err != nil {
		msg := fmt.Sprintf("Error parsing token: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	if token.Valid {
		tokenSubject, err := token.Claims.GetSubject()
		if err != nil {
			msg := "Error getting token subject"
			respondWithError(w, http.StatusBadRequest, msg)
			return
		}
		userId, err := strconv.Atoi(tokenSubject)
		if err != nil {
			msg := fmt.Sprintf("Error parsing userId from tokenSubject :%s", tokenSubject)
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
			wrapper.users[userId-1].PasswordHash = string(passwordHash)
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

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		w.Write(dat)
	} else {
		msg := "Bearer token is invalid"
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}
}

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

	currentTimeUTC := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(currentTimeUTC),
		ExpiresAt: jwt.NewNumericDate(currentTimeUTC.Add(900 * time.Second)),
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
