package handle

import (
	"bytes"
	"errors"

	"github.com/go-playground/validator/v10"
)

func loginValid(l *Login) error {
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
		return errors.New(out.String())
	} else {
		return nil
	}
}

func signupValid(su *SignUp) error {
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
		return errors.New(out.String())
	} else {
		return nil
	}
}

func userPatchValid(up *UserPatch) error {
	validate := validator.New()
	if err := validate.Struct(up); err != nil {
		var out bytes.Buffer
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				switch fe.Field() {
				case "Name":
					out.WriteString("名前は20字以内で入力してください。")
				case "Email":
					if fe.Tag() == "email" {
						out.WriteString("有効なEmailを入力してください。")
					} else if fe.Tag() == "max" {
						out.WriteString("Emailは50字以内で入力してください。")
					}
				}
			}
		}
		return errors.New(out.String())
	} else {
		return nil
	}
}
