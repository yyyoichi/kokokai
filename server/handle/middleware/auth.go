package middleware

import (
	"kokokai/server/auth"
	"kokokai/server/handle"
	"net/http"
)

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		a := auth.NewAuth()
		if err := a.CheckAuth(authHeader); err != nil {
			var res handle.Response
			switch err.Error() {
			case "Unauthorized":
				res = handle.Response{Status: "ログインしてください。"}
			case "invalid_request":
				res = handle.Response{Status: "認証に失敗しました。ログインし直してください。"}
			}
			res.Error(&w)
			return
		}
	})
}
