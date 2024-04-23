package server

import "net/http"

func (wrapper *apiWrapper) resetHandler(w http.ResponseWriter, r *http.Request) {
	wrapper.fileServerHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
