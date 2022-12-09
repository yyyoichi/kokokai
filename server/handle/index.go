package handle

import (
	"encoding/json"
	"kokokai/server/auth"
	cke "kokokai/server/handle/cookie"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
)

func Index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Response{Status: "ok"})
}

type SessionResponse struct {
	Status string `json:"status"`
	UserId string `json:"userId"`
	Name   string `json:"name"`
}

func UserSessionFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	tokenCookie, err := cke.FromUserCookie(r)
	if err != nil || tokenCookie.Value == "" {
		// 認証なし
		NewOkResponse(w)
		return
	}
	jt := auth.NewJwtToken(os.Getenv("SECRET"))
	mc, err := jt.ParseToken(tokenCookie.Value)
	if err != nil {
		// 認証切れ
		NewOkResponse(w)
		return
	}
	res := &SessionResponse{Status: "ok", UserId: mc.Id, Name: mc.Name}
	resJson, err := json.Marshal(res)
	if err != nil {
		NewErrorResponse(err.Error(), w)
		return
	}
	w.Write(resJson)
}
