package handle

import (
	"encoding/json"
	cke "kokokai/server/handle/cookie"
	"net/http"

	"github.com/gorilla/csrf"
)

func NewErrorResponse(status string, w http.ResponseWriter) {
	res := &Response{Status: status}
	json, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(w, string(json), http.StatusBadRequest)
}

func NewOkResponse(w http.ResponseWriter) {
	res := Response{Status: "ok"}
	resJson, err := json.Marshal(res)
	if err != nil {
		NewErrorResponse(err.Error(), w)
		return
	}
	w.Write(resJson)
}

type Response struct {
	Status string `json:"status"`
}

type AuthResponse struct {
	w http.ResponseWriter
	r *http.Request
}

func (ar *AuthResponse) setJWTCookie(token string) {
	c, err := cke.UpdateUserCookie(ar.r, token)
	if err != nil {
		NewErrorResponse(err.Error(), ar.w)
		return
	}
	// jwtをcookieに保存
	http.SetCookie(ar.w, c)
}

func (ar AuthResponse) setCSRFToken() {
	ar.w.Header().Set("X-CSRF-Token", csrf.Token(ar.r))
}

func (ar AuthResponse) writeOk() {
	NewOkResponse(ar.w)
}
