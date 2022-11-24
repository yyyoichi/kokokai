package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtToken struct {
	secret string
}

func NewJwtToken(secret string) *JwtToken {
	return &JwtToken{secret: secret}
}

type myClaims struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func (jt *JwtToken) Generate(u *User) (*string, error) {
	mc := &myClaims{Id: u.Id, Name: u.Name}
	mc.ExpiresAt = jwt.NewNumericDate(time.Now().AddDate(0, 0, 1))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mc)
	tokenString, err := token.SignedString([]byte(jt.secret))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (jt *JwtToken) ParseToken(tokenString string) (*string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jt.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if mc, ok := token.Claims.(*myClaims); ok && token.Valid {
		return &mc.Id, nil
	} else {
		return nil, errors.New("invalid")
	}
}
