package handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kokokai/server/db/user"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestLogin(t *testing.T) {
	loadEnv()
	u := user.User{Id: "yyyoichi", Pass: "pa55W@rd"}
	u.Create()
	defer u.Delete()
	reqBody := bytes.NewBufferString(fmt.Sprintf(`{"id":"%s","pass":"%s"}`, u.Id, u.Pass))
	req := httptest.NewRequest(http.MethodPost, "http://localhost:3000/login", reqBody)

	got := httptest.NewRecorder()
	LoginFunc(got, req)

	var lr LoginResponse
	if err := json.NewDecoder(got.Body).Decode(&lr); err != nil {
		t.Error(err)
	}
	t.Log(lr.Status)
	if lr.Status != "ok" {
		t.Errorf("response excepted ok, but got=%s", lr.Status)
	} else {
		t.Logf("token=%s", lr.Token)
	}
}

func TestLoginEmpty(t *testing.T) {
	loadEnv()
	test := []struct {
		buf            string
		expectedStatus string
	}{
		{
			fmt.Sprintf(`{"id":"%s","pass":"%s"}`, "", "demopass"),
			"id を入力してください。",
		},
		{
			fmt.Sprintf(`{"id":"%s","pass":"%s"}`, "demoid", ""),
			"パスワードを入力してください。",
		},
		{
			fmt.Sprintf(`{"id":"%s","pass":"%s"}`, "", ""),
			"id を入力してください。パスワードを入力してください。",
		},
		{
			fmt.Sprintf(`{"id":"%s","pass":"%s"}`, "zzzyyyxxx", "pass"),
			"idが見つかりません。",
		},
	}
	for _, tt := range test {
		testLoginError(tt.buf, tt.expectedStatus, t)
	}
}

func TestInvalidPassLogin(t *testing.T) {
	loadEnv()
	u := &user.User{Id: "temUser", Pass: "AAA"}
	u.Create()
	defer u.Delete()
	bodyBuf := fmt.Sprintf(`{"id":"%s","pass":"%s"}`, "temUser", "BBB")
	expectedStatus := "パスワードが違います。"
	testLoginError(bodyBuf, expectedStatus, t)
}

func testLoginError(bodyBuf, expectedStatus string, t *testing.T) {
	reqBody := bytes.NewBufferString(bodyBuf)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:3000/login", reqBody)

	got := httptest.NewRecorder()
	LoginFunc(got, req)

	var lr LoginResponse
	if err := json.NewDecoder(got.Body).Decode(&lr); err != nil {
		t.Error(err)
	}
	if lr.Status != expectedStatus {
		t.Errorf("expect err '%s' but got=%s", expectedStatus, lr.Status)
	}
}

func TestSignUp(t *testing.T) {
	loadEnv()
	u := user.User{Id: "o123456789o123456788", Pass: "pa55Ward"}
	defer u.Delete()
	bodyBuf := fmt.Sprintf(`{"id":"%s","pass1":"%s", "pass2":"%s"}`, u.Id, u.Pass, u.Pass)
	testSignUpError(bodyBuf, "ok", t)
}

func testSignUpError(bodyBuf, expectedStatus string, t *testing.T) {
	reqBody := bytes.NewBufferString(bodyBuf)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:3000/signup", reqBody)

	got := httptest.NewRecorder()
	SignUpFunc(got, req)

	var lr LoginResponse
	if err := json.NewDecoder(got.Body).Decode(&lr); err != nil {
		t.Error(err)
	}
	if lr.Status != expectedStatus {
		t.Errorf("response excepted %s, but got=%s", expectedStatus, lr.Status)
	}
}
