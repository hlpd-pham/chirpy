package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hlpd-pham/chirpy/server/auth"
)

func (wrapper *apiWrapper) refresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetToken(r, wrapper.jwtSecret, auth.CHIRPY_REFRESH_ISSUER)
	if err != nil {
		msg := fmt.Sprintf("error parsing refresh token: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	if _, ok := wrapper.revokedTokens[token.Raw]; ok {
		msg := "refresh token is revoked"
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	tokenSubject, err := token.Claims.GetSubject()
	if err != nil {
		msg := fmt.Sprintf("error retrieving subject from token: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	accessToken, err := auth.GetSignedToken(auth.CHIRPY_ACCESS_ISSUER, tokenSubject, wrapper.jwtSecret)
	if err != nil {
		msg := fmt.Sprintf("error signing access token: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	userResponse := refreshResponseBody{
		Token: accessToken,
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
