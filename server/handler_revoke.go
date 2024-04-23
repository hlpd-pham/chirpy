package server

import (
	"fmt"
	"net/http"

	"github.com/hlpd-pham/chirpy/server/auth"
)

func (wrapper *apiWrapper) revoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetToken(r, wrapper.jwtSecret, auth.CHIRPY_REFRESH_ISSUER)
	if err != nil {
		msg := fmt.Sprintf("error parsing refresh token: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	wrapper.revokedTokens[token.Raw] = true

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
}
