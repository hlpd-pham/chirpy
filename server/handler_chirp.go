package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hlpd-pham/chirpy/server/auth"
)

func (wrapper *apiWrapper) getAllChirps(w http.ResponseWriter, _ *http.Request) {
	var chirps []chirp
	for _, value := range wrapper.chirps {
		chirps = append(chirps, value)
	}
	respBody := getAllChirpsResponse{
		Body: chirps,
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
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if _, ok := wrapper.chirps[chirpId]; !ok {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("chirp ID: %d is invalid", chirpId))
		return
	}

	dat, err := json.Marshal(wrapper.chirps[chirpId])
	msg := fmt.Sprintf("Error marshalling JSON response: %s", err)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(dat)
}

func (wrapper *apiWrapper) deleteOneChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if len(id) == 0 {
		respondWithError(w, http.StatusNotFound, "chirp ID is required to delete")
		return
	}

	token, err := auth.GetToken(r, wrapper.jwtSecret, auth.CHIRPY_ACCESS_ISSUER)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	chirpId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userIdToken, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := strconv.Atoi(userIdToken)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, ok := wrapper.chirps[chirpId]
	if !ok {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("chirp ID: %d is invalid", chirpId))
		return
	}

	if userId != chirp.AuthorId {
		respondWithError(w, http.StatusForbidden, "don't delete people's chirps, that's not nice")
		return
	}

	delete(wrapper.chirps, chirpId)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
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

	token, err := auth.GetToken(r, wrapper.jwtSecret, auth.CHIRPY_ACCESS_ISSUER)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	words, err := validateChirp(reqBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	tokenSubject, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := strconv.Atoi(tokenSubject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respBody := chirp{
		Body:     strings.Join(words, " "),
		Id:       wrapper.nextChirpId,
		AuthorId: userId,
	}

	dat, err := json.Marshal(respBody)
	msg := fmt.Sprintf("Error marshalling JSON response: %s", err)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	wrapper.chirps[wrapper.nextChirpId] = respBody
	wrapper.nextChirpId++

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write(dat)
}
