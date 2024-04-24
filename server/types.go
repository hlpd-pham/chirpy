package server

type metricsData struct {
	HitCount int
}

type chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"authorId"`
}

type user struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	IsChirpyRed  bool   `json:"isChirpyRed"`
}

type apiWrapper struct {
	fileServerHits int
	nextChirpId    int
	nextUserId     int
	chirps         map[int]chirp
	users          map[int]user
	jwtSecret      []byte
	polkaKey       []byte
	revokedTokens  map[string]bool
}

type errorResponse struct {
	Message string `json:"message"`
}

type createChirpRequestBody struct {
	Body string `json:"body"`
}

type getAllChirpsResponse struct {
	Body []chirp `json:"body"`
}

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	Email string `json:"email"`
	Id    int    `json:"id"`
}

type loginResponseBody struct {
	Email        string `json:"email"`
	Id           int    `json:"id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type refreshResponseBody struct {
	Token string `json:"token"`
}

type polkaRequestBody struct {
	Event string `json:"event"`
	Data  struct {
		UserId int `json:"user_id"`
	} `json:"data"`
}
