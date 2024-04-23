package server

type metricsData struct {
	HitCount int
}

type chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type user struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

type apiWrapper struct {
	fileServerHits int
	nextChirpId    int
	nextUserId     int
	chirps         []chirp
	users          []user
	jwtSecret      []byte
}

type errorResponse struct {
	Message string `json:"message"`
}

type createChirpRequestBody struct {
	Body string `json:"body"`
}

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	Email string `json:"email"`
	Id    int    `json:"id"`
}

type loginRequestBody struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expiresInSeconds"`
}

type loginResponseBody struct {
	User  userResponse `json:"user"`
	Token string       `json:"token"`
}

type createChirpResponse struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}

type getAllChirpsResponse struct {
	Body []chirp `json:"body"`
}
