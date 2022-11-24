package auth

import (
	"testing"
)

func TestJWTToken(t *testing.T) {
	secret := "secret!!00"
	j := NewJwtToken(secret)
	u := &User{Id: "xxxyyyzzz", Name: "yyyoichi", Email: "yyyoichi@example.com"}
	tokenString, err := j.Generate(u)
	if err != nil {
		t.Error(err)
	}
	id, err := j.ParseToken(*tokenString)
	if err != nil {
		t.Error(err)
	}
	if *id != u.Id {
		t.Errorf("expected: %s, got=%s", u.Id, *id)
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
	_, err = j.ParseToken(*tokenString)
	if err == nil {
		t.Error("??")
	} else {
		t.Log(err)
	}
}
