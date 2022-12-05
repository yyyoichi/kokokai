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

type MyClaims struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func (jt *JwtToken) Generate(id, name string) (*string, error) {
	mc := &MyClaims{Id: id, Name: name}
	mc.ExpiresAt = jwt.NewNumericDate(time.Now().AddDate(0, 0, 1))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mc)
	tokenString, err := token.SignedString([]byte(jt.secret))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (jt *JwtToken) ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jt.secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid")
	}

	if mc, ok := token.Claims.(*MyClaims); ok {
		return mc, nil
	} else {
		return nil, errors.New("invalid")
	}
}
