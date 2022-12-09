package handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kokokai/server/auth"
	"kokokai/server/db/user"
	ctx "kokokai/server/handle/context"
	cke "kokokai/server/handle/cookie"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
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

	var r Response
	if err := json.NewDecoder(got.Body).Decode(&r); err != nil {
		t.Error(err)
	}
	t.Log(r.Status)
	if r.Status != "ok" {
		t.Errorf("response excepted ok, but got=%s", r.Status)
	}
}

var loginTestUnit = []struct {
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

func TestLoginEmpty(t *testing.T) {
	loadEnv()
	test := loginTestUnit
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

	var r Response
	if err := json.NewDecoder(got.Body).Decode(&r); err != nil {
		t.Error(err)
	}
	if r.Status != expectedStatus {
		t.Errorf("expect err '%s' but got=%s", expectedStatus, r.Status)
	}
}

func TestSignUp(t *testing.T) {
	loadEnv()
	u := user.User{Id: "o123456789o123456788", Pass: "pa55Ward"}
	defer u.Delete()
	bodyBuf := fmt.Sprintf(`{"id":"%s","pass1":"%s", "pass2":"%s"}`, u.Id, u.Pass, u.Pass)
	testSignUpError(bodyBuf, "ok", t)
}

var (
	normId        = "o123456789o123456789"
	normPass      = "pa55w0rd"
	sinupTestUnit = []struct {
		buf            string
		expectedStatus string
	}{
		{
			fmt.Sprintf(`{"id":"%s","pass1":"%s","pass2":"%s"}`, "abc3", normPass, normPass),
			"idは20字で入力してください。",
		},
		{
			fmt.Sprintf(`{"id":"%s","pass1":"%s","pass2":"%s"}`, "", normPass, normPass),
			"id を入力してください。",
		},
		{
			fmt.Sprintf(`{"id":"%s","pass1":"%s","pass2":"%s"}`, normId, "normPasspass", "normPasspass"),
			"パスワードは英数字である必要があります。",
		},
		{
			fmt.Sprintf(`{"id":"%s","pass1":"%s","pass2":"%s"}`, normId, "012345678", "012345678"),
			"パスワードは英数字である必要があります。",
		},
		{
			fmt.Sprintf(`{"id":"%s","pass1":"%s","pass2":"%s"}`, normId, "pass1", "pass1"),
			"パスワードは8~24字である必要があります。",
		},
		{
			fmt.Sprintf(`{"id":"%s","pass1":"%s","pass2":"%s"}`, normId, "normPassnormPassnormPass123456", "normPassnormPassnormPass123456"),
			"パスワードは8~24字である必要があります。",
		},
		{
			fmt.Sprintf(`{"id":"%s","pass1":"%s","pass2":"%s"}`, normId, "", ""),
			"パスワードを入力してください。確認用のパスワードを入力してください。",
		},
		{
			fmt.Sprintf(`{"id":"%s","pass1":"%s","pass2":"%s"}`, normId, normPass, "normPass"),
			"パスワードが一致しません。",
		},
		{
			fmt.Sprintf(`{"id":"%s","pass1":"%s","pass2":"%s"}`, normId, normPass, ""),
			"確認用のパスワードを入力してください。",
		},
	}
)

func TestSignUpError(t *testing.T) {
	loadEnv()
	test := sinupTestUnit
	for _, tt := range test {
		testSignUpError(tt.buf, tt.expectedStatus, t)
	}
}

func testSignUpError(bodyBuf, expectedStatus string, t *testing.T) {
	reqBody := bytes.NewBufferString(bodyBuf)
	req := httptest.NewRequest(http.MethodPost, "http://localhost:3000/signup", reqBody)

	got := httptest.NewRecorder()
	SignUpFunc(got, req)

	var r Response
	if err := json.NewDecoder(got.Body).Decode(&r); err != nil {
		t.Error(err)
	}
	if r.Status != expectedStatus {
		t.Errorf("response excepted %s, but got=%s", expectedStatus, r.Status)
	}
}

var userPatchTestUnit = []struct {
	buf               string
	expectedStatus    string
	expectedName      string
	expectedEmail     string
	expectTokenUpdate bool
}{
	{ //正常系
		fmt.Sprintf(`{"name":"%s","email":"%s"}`, "yyyoichi", "yyyoichi@example.com"),
		"ok",
		"yyyoichi",
		"yyyoichi@example.com",
		true,
	},
	{ // validation
		fmt.Sprintf(`{"name":"%s","email":"%s"}`, "abcdefghijabcdefghijdd", "yyyoichi@example.com"),
		"名前は20字以内で入力してください。",
		"yyyoichi",
		"yyyoichi@example.com",
		false,
	},
	{ // validation
		fmt.Sprintf(`{"name":"%s","email":"%s"}`, "yyyoichi", "yyyoichiexample.com"),
		"有効なEmailを入力してください。",
		"yyyoichi",
		"yyyoichi@example.com",
		false,
	},
	{ // validation
		fmt.Sprintf(`{"name":"%s","email":"%s"}`, "yyyoichi", "012345678901234567890123456789012345678900123456789@example.com"),
		"Emailは50字以内で入力してください。",
		"yyyoichi",
		"yyyoichi@example.com",
		false,
	},
	{ // 正常系一部
		fmt.Sprintf(`{"name":"%s"}`, "hogehoge"),
		"ok",
		"hogehoge",
		"yyyoichi@example.com",
		true,
	},
	{ // 正常系一部
		fmt.Sprintf(`{"email":"%s"}`, "hogehoge@example.com"),
		"ok",
		"hogehoge",
		"hogehoge@example.com",
		false,
	},
	{ // 正常系空欄
		fmt.Sprintf(`{"name":"%s", "email":"%s"}`, "", ""),
		"ok",
		"hogehoge",
		"hogehoge@example.com",
		false,
	},
}

func TestUserPath(t *testing.T) {
	loadEnv()
	u := &user.User{Id: "tmpuser", Pass: "pa55w0rd", Name: "", Email: "example@example.com"}
	if err := u.Create(); err != nil {
		t.Error(err)
	}
	defer u.Delete()

	r := mux.NewRouter()
	r.HandleFunc("/users/{userId}", UserFunc)
	s := os.Getenv("SECRET")
	jt := auth.NewJwtToken(s)
	token, _ := jt.Generate(u.Id, "")
	beforeClaims, _ := jt.ParseToken(*token)
	for i, tt := range userPatchTestUnit {
		reqBody := bytes.NewBufferString(tt.buf)
		req := httptest.NewRequest(http.MethodPatch, `/users/`+u.Id, reqBody)

		// context にアップデート前のユーザ情報を仕込む
		uctx := ctx.NewUserContext(req.Context(), beforeClaims)
		got := httptest.NewRecorder()
		// Cookkieを仕込む
		c := cke.NewUserCookie(*token)
		req.AddCookie(c)

		r.ServeHTTP(got, req.WithContext(uctx))
		if err := u.GetById(); err != nil {
			t.Errorf("%d: %s", i, err)
		}
		var res Response
		if err := json.NewDecoder(got.Body).Decode(&res); err != nil {
			t.Error(err)
		}
		if res.Status != tt.expectedStatus {
			t.Errorf("%d: excepted status %s but got=%s", i, tt.expectedStatus, res.Status)
		}
		if u.Name != tt.expectedName {
			t.Errorf("%d: excepted name %s but got=%s", i, tt.expectedName, u.Name)
		}
		if u.Email != tt.expectedEmail {
			t.Errorf("%d: excepted email %s but got=%s", i, tt.expectedEmail, u.Email)
		}
		if tt.expectTokenUpdate {
			cs := got.Result().Cookies()
			mc, err := jt.ParseToken(cs[0].Value)
			if err != nil {
				t.Error(err)
			}
			if mc.Name != tt.expectedName {
				t.Errorf("%d: expected name in token %s but got=%s", i, tt.expectedName, mc.Name)
			}
		}
	}
}

func TestSession(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	got := httptest.NewRecorder()
	s := []byte(os.Getenv("CSRF_SECRET"))
	csrfMiddleware := csrf.Protect(s, csrf.Secure(false))
	fu := csrfMiddleware(http.HandlerFunc(UserSessionFunc))
	fu.ServeHTTP(got, req)
	var res SessionResponse
	if err := json.NewDecoder(got.Body).Decode(&res); err != nil {
		t.Error(err)
	}
	if res.Status != "ok" {
		t.Error("no ok")
	}
	t.Log(got.Result().Header.Get("X-CSRF-Token"))
}
