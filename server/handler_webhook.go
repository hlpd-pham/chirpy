package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hlpd-pham/chirpy/server/auth"
)

func (wrapper *apiWrapper) polka(w http.ResponseWriter, r *http.Request) {
	err := auth.ValidatePolkaKey(r, wrapper.polkaKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := polkaRequestBody{}
	err = decoder.Decode(&reqBody)
	if err != nil {
		msg := fmt.Sprintf("Error decoding request body: %s", err)
		respondWithError(w, http.StatusServiceUnavailable, msg)
		return
	}

	if reqBody.Event == "user.upgraded" && reqBody.Data.UserId > 0 {
		if user, ok := wrapper.users[reqBody.Data.UserId]; ok {
			user.IsChirpyRed = true
			wrapper.users[reqBody.Data.UserId] = user
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "application/json")
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Add("Content-Type", "application/json")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
}
