package middleware

import (
	"kokokai/server/auth"
	"kokokai/server/handle"
	ctx "kokokai/server/handle/context"
	cke "kokokai/server/handle/cookie"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := cke.FromUserCookie(r)
		if err != nil || tokenCookie.Value == "" {
			st := "ログインしてください。"
			handle.NewErrorResponse(st, w)
			return
		}
		secret := os.Getenv("SECRET")
		j := auth.NewJwtToken(secret)
		mc, err := j.ParseToken(tokenCookie.Value)
		if err != nil {
			var st string
			switch err.Error() {
			case "unexpected signing method":
				st = "不正アクセス"
			case "invalid":
				st = "認証に失敗しました。ログインし直してください。"
			default:
				st = "予期せぬエラーが発生しました。"
			}
			handle.NewErrorResponse(st, w)
			return
		}
		// context にユーザ情報格納
		userCxt := ctx.NewUserContext(r.Context(), mc)
		// 認証成功
		vars := mux.Vars(r)
		if vars["userId"] != "" {
			// ログインユーザとリクエスト対象のユーザが一致しない。
			if vars["userId"] != mc.Id {
				st := "不正な操作です。ログインし直してください。"
				handle.NewErrorResponse(st, w)
				return
			}
		}
		next.ServeHTTP(w, r.WithContext(userCxt))
	})
}
