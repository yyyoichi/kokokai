package middleware

import (
	"kokokai/server/auth"
	"kokokai/server/handle"
	"net/http"

	"github.com/gorilla/mux"
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
		// 認証成功
		vars := mux.Vars(r)
		if vars["userId"] != "" {
			// ログインユーザとリクエスト対象のユーザが一致しない。
			if vars["userId"] != a.UserId {
				res := handle.Response{Status: "不正な操作です。ログインし直してください。"}
				res.Error(&w)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
