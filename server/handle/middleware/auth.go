package middleware

import (
	"kokokai/server/auth"
	"kokokai/server/handle"
	ctx "kokokai/server/handle/context"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Brearer ") {
			res := handle.Response{Status: "ログインしてください。"}
			res.Error(&w)
			return
		}
		secret := os.Getenv("SECRET")
		j := auth.NewJwtToken(secret)
		mc, err := j.ParseToken(authHeader[7:])
		if err != nil {
			var res handle.Response
			switch err.Error() {
			case "unexpected signing method":
				res = handle.Response{Status: "不正アクセス"}
			case "invalid":
				res = handle.Response{Status: "認証に失敗しました。ログインし直してください。"}
			default:
				res = handle.Response{Status: "予期せぬエラーが発生しました。"}
			}
			res.Error(&w)
			return
		}
		// context にユーザ情報格納
		userCxt := ctx.NewUserContext(r.Context(), mc)
		// 認証成功
		vars := mux.Vars(r)
		if vars["userId"] != "" {
			// ログインユーザとリクエスト対象のユーザが一致しない。
			if vars["userId"] != mc.Id {
				res := handle.Response{Status: "不正な操作です。ログインし直してください。"}
				res.Error(&w)
				return
			}
		}
		next.ServeHTTP(w, r.WithContext(userCxt))
	})
}
