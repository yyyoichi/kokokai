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
		reqBody := bytes.NewBufferString(tt.buf)
		req := httptest.NewRequest(http.MethodPost, "http://localhost:3000/login", reqBody)

		got := httptest.NewRecorder()
		LoginFunc(got, req)

		var lr LoginResponse
		if err := json.NewDecoder(got.Body).Decode(&lr); err != nil {
			t.Error(err)
		}
		if lr.Status != tt.expectedStatus {
			t.Errorf("expect err '%s' but got=%s", tt.expectedStatus, lr.Status)
		}
	}
}
