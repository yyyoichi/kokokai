package handle

import (
	"bytes"
	"encoding/json"
	"errors"
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
		secret := os.Getenv("SECRET")
		j := auth.NewJwtToken(secret)
		tokenString, err := j.Generate(user.Id, user.Name)
		if err != nil {
			res := Response{err.Error()}
			res.resError(&w)
			return
		}
		res := &LoginResponse{"ok", *tokenString}
		json, err := json.Marshal(res)
		if err != nil {
			res := Response{err.Error()}
			res.resError(&w)
			return
		}
		w.Write(json)
	default:
		res := Response{"permits only POST"}
		res.resError(&w)
		return
	}
}
