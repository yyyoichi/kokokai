package handle

import (
	"bytes"
	"errors"

	"github.com/go-playground/validator/v10"
)

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
