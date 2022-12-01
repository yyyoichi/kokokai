package handle

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"kokokai/server/db/user"
	"net/http"

	"github.com/go-playground/validator/v10"
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
			http.Error(w, fmt.Sprintf(`{"status":"%s"}`, err), http.StatusInternalServerError)
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
			http.Error(w, fmt.Sprintf(`{"status:":"%s"}`, out.String()), http.StatusBadRequest)
			return
		}
		// 入力値正常
		user := &user.User{Id: l.Id, Pass: l.Pass}
		if err := user.GetByPass(); err != nil {
			switch err.Error() {
			case "wrong-pass":
				http.Error(w, fmt.Sprintf(`{"status:":"%s"}`, "パスワードが違います。"), http.StatusBadRequest)
				return
			case "no-email":
				http.Error(w, fmt.Sprintf(`{"status:":"%s"}`, "Id。"), http.StatusBadRequest)
				return
			}
		}

	}
}
