package auth

import (
	"testing"
)

func TestJWTToken(t *testing.T) {
	secret := "secret!!00"
	j := NewJwtToken(secret)
	tokenString, err := j.Generate(&User{Name: "yyyoichi", Email: "yyyoichi@example.com"})
	if err != nil {
		t.Error(err)
	}
	t.Log(*tokenString)
	err = j.ParseToken(*tokenString)
	if err != nil {
		t.Error(err)
	}
}

func TestInvalidToken(t *testing.T) {
	secret := "secret!!00"
	j := NewJwtToken(secret)
	tokenString, err := j.Generate(&User{Name: "yyyoichi", Email: "yyyoichi@example.com"})
	if err != nil {
		t.Error(err)
	}
	t.Log(*tokenString)
	j.secret = "secret!!11"
	err = j.ParseToken(*tokenString)
	if err == nil {
		t.Error("??")
	} else {
		t.Log(err)
	}
}
