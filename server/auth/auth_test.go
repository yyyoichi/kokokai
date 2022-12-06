package auth

import (
	"os"
	"testing"
)

func TestAuthCheck(t *testing.T) {
	secret := os.Getenv("SECRET")
	jt := NewJwtToken(secret)
	tokenString, err := jt.Generate("test1", "yyyoichi")
	if err != nil {
		t.Error(err)
	}
	at := NewAuth()
	err = at.CheckAuth("Bearer " + *tokenString)
	if err != nil {
		t.Error(err)
	}
	t.Log(at.UserId)
}

func TestEmpTokenCheck(t *testing.T) {
	at := NewAuth()
	err := at.CheckAuth("")
	if err == nil {
		t.Error("not err")
	}
	if err.Error() != "Unauthorized" {
		t.Errorf("expected Unauthorized but got=%s", err)
	}
}

func TestInValidTokenCheck(t *testing.T) {
	at := NewAuth()
	err := at.CheckAuth("Bearer ")
	if err == nil {
		t.Error("not err")
	}
	if err.Error() != "invalid_request" {
		t.Errorf("expected Unauthorized but got=%s", err)
	}
}
