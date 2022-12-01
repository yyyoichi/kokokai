package handle

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"kokokai/server/auth"
	"kokokai/server/db/user"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
)

type Login struct {
	Id   string `validate:"required"`
	Pass string `validate:"required"`
}

type LoginResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

func (lr *LoginResponse) resWithJWT(w http.ResponseWriter, user *user.User) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// jwt作成
	secret := os.Getenv("SECRET")
	j := auth.NewJwtToken(secret)
	tokenString, err := j.Generate(user.Id, user.Name)
	if err != nil {
		res := Response{err.Error()}
		res.resError(&w)
		return
	}
	lr.Token = *tokenString
	resJson, err := json.Marshal(lr)
	if err != nil {
		res := Response{err.Error()}
		res.resError(&w)
		return
	}
	w.Write(resJson)
}

func LoginFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch r.Method {
	case http.MethodPost:
		var l Login
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			res := &Response{"need id and pass field"}
			res.resError(&w)
			return
		}

		validate := validator.New()
		if err := validate.Struct(l); err != nil {
			var out bytes.Buffer
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				for _, fe := range ve {
					switch fe.Field() {
					case "Id":
						out.WriteString("id を入力してください。")
					case "Pass":
						out.WriteString("パスワードを入力してください。")
					}
				}
			}
			res := &Response{out.String()}
			res.resError(&w)
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
			res := Response{out.String()}
			res.resError(&w)
			return
		}
		// DBから取得正常
		// jwt作成
		res := LoginResponse{Status: "ok"}
		res.resWithJWT(w, user)
	default:
		res := Response{"permits only POST"}
		res.resError(&w)
		return
	}
}

type SignUp struct {
	Id    string `validate:"required,len=20"`
	Pass1 string `validate:"required,alphanum,min=8,max=24"`
	Pass2 string `validate:"required,eqfield=Pass1"`
}

func SignUpFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch r.Method {
	case http.MethodPost:
		var su SignUp
		if err := json.NewDecoder(r.Body).Decode(&su); err != nil {
			res := &Response{"need id and pass field"}
			res.resError(&w)
			return
		}

		validate := validator.New()
		if err := validate.Struct(su); err != nil {
			var out bytes.Buffer
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				for _, fe := range ve {
					switch fe.Field() {
					case "Id":
						if fe.Tag() == "len" {
							out.WriteString("idは20字で入力してください。")
						} else {
							out.WriteString("id を入力してください。")
						}
					case "Pass1":
						if fe.Tag() == "alphanum" {
							out.WriteString("id は英数字である必要があります。")
						} else if fe.Tag() == "min" || fe.Tag() == "max" {
							out.WriteString("id は8~24字である必要があります。")
						} else {
							out.WriteString("パスワードを入力してください。")
						}
					case "Pass2":
						if fe.Tag() == "eqfield" {
							out.WriteString("パスワードが一致しません。")
						} else {
							out.WriteString("確認用のパスワードを入力してください。")
						}
					}
				}
			}
			res := &Response{out.String()}
			res.resError(&w)
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
			res := Response{out.String()}
			res.resError(&w)
			return
		}
		// DBに新しいユーザを作成完了
		// jwt作成
		res := LoginResponse{Status: "ok"}
		res.resWithJWT(w, user)
	default:
		res := Response{"permits only POST"}
		res.resError(&w)
		return
	}
}
