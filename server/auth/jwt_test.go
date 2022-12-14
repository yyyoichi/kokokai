package auth

import (
	"testing"
)

func TestJWTToken(t *testing.T) {
	secret := "secret!!00"
	j := NewJwtToken(secret)
	tokenString, err := j.Generate("xxxyyyzzz", "yyyoichi")
	if err != nil {
		t.Error(err)
	}
	mc, err := j.ParseToken(*tokenString)
	if err != nil {
		t.Error(err)
	}
	if mc.Id != "xxxyyyzzz" {
		t.Errorf("expected: 'xxxyyyzzz', got=%s", mc.Id)
	}
}

func TestInvalidToken(t *testing.T) {
	secret := "secret!!00"
	j := NewJwtToken(secret)
	tokenString, err := j.Generate("xxxyyyzzz", "yyyoichi")
	if err != nil {
		t.Error(err)
	}
	t.Log(*tokenString)
	j.secret = "secret!!11"
	_, err = j.ParseToken(*tokenString)
	if err == nil {
		t.Error("??")
	} else {
		t.Log(err)
	}
}
