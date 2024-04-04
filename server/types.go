package server

type metricsData struct {
	HitCount int
}

type chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type user struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type apiWrapper struct {
	fileServerHits int
	nextChirpId    int
	nextUserId     int
	chirps         []chirp
	users          []user
}

type errorResponse struct {
	Error error `json:"error"`
}

type createChirpRequestBody struct {
	Body string `json:"body"`
}

type createUserRequestBody struct {
	Email string `json:"email"`
}

type createChirpResponse struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}

type getAllChirpsResponse struct {
	Body []chirp `json:"body"`
}
