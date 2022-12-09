package middleware

import (
	"encoding/json"
	"fmt"
	"kokokai/server/auth"
	"kokokai/server/db/user"
	"kokokai/server/handle"
	ctx "kokokai/server/handle/context"
	cke "kokokai/server/handle/cookie"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

type authTest struct {
	req *http.Request
	got *httptest.ResponseRecorder
}

func (at *authTest) addAuthCookie(token string) {
	c := cke.NewUserCookie(token)
	at.req.AddCookie(c)
}

func NewAuthTest(urlUserId string) *authTest {
	at := &authTest{}
	at.req = httptest.NewRequest(http.MethodPost, "/users/"+urlUserId, nil)
	at.got = httptest.NewRecorder()
	return at
}

func useRouter() *mux.Router {
	loadEnv()
	r := mux.NewRouter()
	r.Use(MiddlewareAuth)
	h := func(w http.ResponseWriter, r *http.Request) {
		mc, ok := ctx.FromUserContext(r.Context())
		if ok {
			fmt.Println(mc.Id + "," + mc.Name)
		}
		res := handle.Response{Status: "ok"}
		resJson, err := json.Marshal(res)
		if err != nil {
			res := handle.Response{Status: err.Error()}
			res.Error(&w)
			return
		}
		w.Write(resJson)
	}
	r.HandleFunc("/users/{userId}", h)
	r.HandleFunc("/users/", h)
	http.Handle("/", r)
	return r
}

var router *mux.Router = useRouter()

func TestAuthMiddleware(t *testing.T) {
	u := user.User{Id: "yyyoichi"}
	s := os.Getenv("SECRET")
	auth := auth.NewJwtToken(s)
	next := NewAuthTest(u.Id)

	token, err := auth.Generate(u.Id, u.Name)
	if err != nil {
		t.Error(err)
	}
	next.addAuthCookie(*token)

	router.ServeHTTP(next.got, next.req)
	var res handle.Response
	if err := json.NewDecoder(next.got.Body).Decode(&res); err != nil {
		t.Error(err)
	}
	if res.Status != "ok" {
		t.Errorf("expected ok. but got=%s", res.Status)
	}
}

func TestMiddlewareAuth(t *testing.T) {
	id := "yyyoichi"
	s := os.Getenv("SECRET")
	test := []struct {
		pathId         string
		tokenId        string
		secret         string
		hasHead        bool
		expectedStatus string
	}{
		{id, id, s, true, "ok"},           // 正常系
		{"", id, s, true, "ok"},           // 正常系パラメータなし
		{id, id, s, false, "ログインしてください。"}, //認証ヘッダなし
		{id, id, "hoge", true, "認証に失敗しました。ログインし直してください。"}, //不正jwt.署名違い
		{"other", id, s, true, "不正な操作です。ログインし直してください。"},   //別User操作
	}
	for _, tt := range test {
		next := NewAuthTest(tt.pathId)
		if tt.hasHead {
			a := auth.NewJwtToken(tt.secret)
			token, err := a.Generate(tt.tokenId, "")
			if err != nil {
				t.Error(err)
			}
			next.addAuthCookie(*token)
		}

		router.ServeHTTP(next.got, next.req)
		var res handle.Response
		if err := json.NewDecoder(next.got.Body).Decode(&res); err != nil {
			t.Error(err)
		}
		if res.Status != tt.expectedStatus {
			t.Errorf("expected status '%s'. but got=%s", tt.expectedStatus, res.Status)
		}
	}
}
