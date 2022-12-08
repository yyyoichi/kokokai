package handle

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"kokokai/server/auth"
	"kokokai/server/db/user"
	ctx "kokokai/server/handle/context"
	cke "kokokai/server/handle/cookie"
	sess "kokokai/server/handle/session"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Login struct {
	Id   string `validate:"required"`
	Pass string `validate:"required"`
}

type LoginResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

func (lr *LoginResponse) resWithJWT(w http.ResponseWriter, r *http.Request, user *user.User) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// jwt作成
	secret := os.Getenv("SECRET")
	j := auth.NewJwtToken(secret)
	tokenString, err := j.Generate(user.Id, user.Name)
	if err != nil {
		res := Response{err.Error()}
		res.Error(&w)
		return
	}
	// jwtをcookieに保存
	c := cke.NewUserCookie(*tokenString)
	http.SetCookie(w, c)
	// csrfTokenをセッションに保存
	s := sess.NewUserCSRFToken(r)
	s.Save(r, w)
	// body に返却
	lr.Token = s.Values["csrftoken"].(string)
	resJson, err := json.Marshal(lr)
	if err != nil {
		res := Response{err.Error()}
		res.Error(&w)
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
			res.Error(&w)
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
			res.Error(&w)
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
			res.Error(&w)
			return
		}
		// DBから取得正常
		// jwt作成
		res := LoginResponse{Status: "ok"}
		res.resWithJWT(w, r, user)
	default:
		res := Response{"permits only POST"}
		res.Error(&w)
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
			res := &Response{"need id and pass field"}
			res.Error(&w)
			return
		}

		validate := validator.New()
		validate.RegisterValidation("alphanumary", customAlphanumary)
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
						if fe.Tag() == "alphanumary" {
							out.WriteString("パスワードは英数字である必要があります。")
						} else if fe.Tag() == "min" || fe.Tag() == "max" {
							out.WriteString("パスワードは8~24字である必要があります。")
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
			res.Error(&w)
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
			res.Error(&w)
			return
		}
		// DBに新しいユーザを作成完了
		// jwt作成
		res := LoginResponse{Status: "ok"}
		res.resWithJWT(w, r, user)
	default:
		res := Response{"permits only POST"}
		res.Error(&w)
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
			res := &Response{"Need name or email field"}
			res.Error(&w)
			return
		}

		if err := userPatchValid(&up); err != nil {
			res := &Response{Status: err.Error()}
			res.Error(&w)
			return
		}
		// バリデーションチェック完了。入力正常。
		// ユーザ情報アップデート
		vars := mux.Vars(r)
		userId := vars["userId"]
		u := &user.User{Id: userId}
		updateColumn := up.getPatchColumns(u)
		if err := u.Update(updateColumn); err != nil {
			res := Response{err.Error()}
			res.Error(&w)
			return
		}

		var tokenString string = ""
		if up.Name != "" {
			// トークン情報をアップデートする
			mc, ok := ctx.FromUserContext(r.Context())
			if !ok {
				res := Response{Status: "予期せぬエラーが発生しました。"}
				res.Error(&w)
				return
			}
			secret := os.Getenv("SECRET")
			jt := auth.NewJwtToken(secret)
			token, err := jt.UpdateName(mc, u.Name)
			if err != nil {
				if !ok {
					res := Response{Status: err.Error()}
					res.Error(&w)
					return
				}
			}
			c, err := cke.UpdateUserCookie(r, *token)
			if err != nil {
				res := Response{Status: err.Error()}
				res.Error(&w)
				return
			}
			http.SetCookie(w, c)
			tokenString = *token
		}

		res := LoginResponse{Status: "ok", Token: tokenString}
		resJson, err := json.Marshal(res)
		if err != nil {
			res := Response{err.Error()}
			res.Error(&w)
			return
		}
		w.Write(resJson)
	default:
		res := Response{"permits only MethodPatch"}
		res.Error(&w)
		return
	}
}
