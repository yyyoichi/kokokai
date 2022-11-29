package auth

import (
	"errors"
	"os"
)

type Auth struct {
	UserId string
}

func NewAuth() *Auth {
	return &Auth{}
}

func (at *Auth) CheckAuth(authorizationHeader string) error {
	// authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return errors.New("Unauthorized")
	}
	secret := os.Getenv("SECRET")
	j := NewJwtToken(secret)
	userId, err := j.ParseToken(authorizationHeader[7:])
	if err != nil {
		return errors.New("invalid_request")
	}
	at.UserId = *userId
	return nil
}
