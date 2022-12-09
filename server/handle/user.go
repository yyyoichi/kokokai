package handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kokokai/server/auth"
	"kokokai/server/db/user"
	ctx "kokokai/server/handle/context"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Login struct {
	Id   string `validate:"required"`
	Pass string `validate:"required"`
}

func LoginFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch r.Method {
	case http.MethodPost:
		var l Login
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			NewErrorResponse("need id and pass field", w)
			return
		}

		// バリデーションチェック
		if err := loginValid(&l); err != nil {
			NewErrorResponse(err.Error(), w)
			return
		}
		// 入力値正常
		user := &user.User{Id: l.Id, Pass: l.Pass}
		if err := user.GetByPass(); err != nil {
			var out bytes.Buffer
			switch err.Error() {
			case "wrong-pass":
				out.WriteString("パスワードが違います。")
			case "no-id":
				out.WriteString("idが見つかりません。")
			default:
				out.WriteString(err.Error())
			}
			NewErrorResponse(out.String(), w)
			return
		}
		// DBから取得正常
		// jwt作成
		secret := os.Getenv("SECRET")
		jt := auth.NewJwtToken(secret)
		token, err := jt.Generate(user.Id, user.Name)
		if err != nil {
			NewErrorResponse("予期せぬエラーが発生しました: "+err.Error(), w)
			return
		}
		res := AuthResponse{r: r, w: w}
		res.setJWTCookie(*token)
		res.setCSRFToken()
		res.writeOk()
	default:
		NewErrorResponse("permits only MethodPOST", w)
		return
	}
}

type SignUp struct {
	Id    string `validate:"required,len=20"`
	Pass1 string `validate:"required,alphanumary,min=8,max=24"`
	Pass2 string `validate:"required,eqfield=Pass1"`
}

func SignUpFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch r.Method {
	case http.MethodPost:
		var su SignUp
		if err := json.NewDecoder(r.Body).Decode(&su); err != nil {
			NewErrorResponse("need id, pass1 and pass2 field", w)
			return
		}

		if err := signupValid(&su); err != nil {
			NewErrorResponse(err.Error(), w)
			return
		}
		// バリデーションチェック完了。入力正常。
		// ユーザ作成
		user := &user.User{Id: su.Id, Pass: su.Pass1}
		if err := user.Create(); err != nil {
			var out bytes.Buffer
			switch err.Error() {
			case fmt.Sprintf("%s is exists", user.Id):
				out.WriteString("すでにidが存在しています。")
			default:
				out.WriteString("予期せぬエラーが発生しました。")
			}
			NewErrorResponse(out.String(), w)
			return
		}
		// DBに新しいユーザを作成完了
		// jwt作成
		secret := os.Getenv("SECRET")
		jt := auth.NewJwtToken(secret)
		token, err := jt.Generate(user.Id, user.Name)
		if err != nil {
			NewErrorResponse("予期せぬエラーが発生しました: "+err.Error(), w)
			return
		}
		res := AuthResponse{r: r, w: w}
		res.setJWTCookie(*token)
		res.setCSRFToken()
		res.writeOk()
	default:
		NewErrorResponse("permits only MethodPOST", w)
		return
	}
}

type UserPatch struct {
	Name  string `validate:"omitempty,max=20"`
	Email string `validate:"omitempty,email,max=50"`
}

func (up *UserPatch) getPatchColumns(u *user.User) []user.ColumnName {
	c := make([]user.ColumnName, 0)
	if up.Name != "" {
		u.Name = up.Name
		c = append(c, user.ColumnName("name"))
	}
	if up.Email != "" {
		u.Email = up.Email
		c = append(c, user.ColumnName("email"))
	}
	return c
}

func UserFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch r.Method {
	case http.MethodPatch:
		var up UserPatch
		if err := json.NewDecoder(r.Body).Decode(&up); err != nil {
			NewErrorResponse("Need name or email field", w)
			return
		}

		if err := userPatchValid(&up); err != nil {
			NewErrorResponse(err.Error(), w)
			return
		}
		// バリデーションチェック完了。入力正常。
		// ユーザ情報アップデート
		vars := mux.Vars(r)
		userId := vars["userId"]
		u := &user.User{Id: userId}
		updateColumn := up.getPatchColumns(u)
		if err := u.Update(updateColumn); err != nil {
			NewErrorResponse(err.Error(), w)
			return
		}
		// jwtのアップデートが不要
		if up.Name == "" {
			NewOkResponse(w)
			return
		}
		// トークン情報をアップデートする
		// 前のJWTの中身を取り出す
		mc, ok := ctx.FromUserContext(r.Context())
		if !ok {
			NewErrorResponse("予期せぬエラーが発生しました。", w)
			return
		}
		secret := os.Getenv("SECRET")
		jt := auth.NewJwtToken(secret)
		// 前のトークン情報の名前部分を書き換える
		token, err := jt.UpdateName(mc, u.Name)
		if err != nil {
			NewErrorResponse(err.Error(), w)
			return
		}
		res := AuthResponse{r: r, w: w}
		res.updateJWTCookie(*token)
		res.setCSRFToken()
		res.writeOk()
	default:
		NewErrorResponse("permits only MethodPatch", w)
		return
	}
}

func UserSessionFunc(w http.ResponseWriter, r *http.Request) {

}
