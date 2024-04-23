package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	CHIRPY_ACCESS_ISSUER  = "chirpy-access"
	CHIRPY_REFRESH_ISSUER = "chirpy-refresh"
)

func GetSignedToken(issuer, subject string, signKey []byte) (string, error) {
	var tokenDuration time.Duration
	switch issuer {
	case CHIRPY_ACCESS_ISSUER:
		tokenDuration = time.Hour
	case CHIRPY_REFRESH_ISSUER:
		tokenDuration = time.Hour * 24 * 60
	default:
		return "", fmt.Errorf("invalid issuer: %s", issuer)
	}

	currentTimeUTC := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(currentTimeUTC),
		ExpiresAt: jwt.NewNumericDate(currentTimeUTC.Add(tokenDuration)),
		Subject:   subject,
	})
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GetToken(r *http.Request, signKey []byte, tokenType string) (*jwt.Token, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("could not find auth header from request")
	}
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, fmt.Errorf("auth header is not formatted correctly")
	}

	token, err := jwt.ParseWithClaims(
		headerParts[1],
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return signKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %s", err)
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return nil, fmt.Errorf("error getting token issuer")
	}
	if issuer != tokenType {
		return nil, fmt.Errorf("token is %s token, %s token is required", issuer, tokenType)
	}
	if !token.Valid {
		return nil, fmt.Errorf("bearer token is invalid")
	}
	return token, nil
}
