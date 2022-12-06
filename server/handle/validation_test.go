package handle

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestValidation(t *testing.T) {
	// myEmail := ""
	validate := validator.New()
	var ve validator.ValidationErrors
	if err := validate.Var(nil, "omitempty,email"); err != nil {
		errors.As(err, &ve)
		t.Error(ve)
		for _, fe := range ve {
			t.Error(len(fe.Field()))
		}
	}
}

func TestUserPatch(t *testing.T) {
	test := userPatchTestUnit
	for i, tt := range test {
		var up UserPatch
		if err := json.NewDecoder(strings.NewReader(tt.buf)).Decode(&up); err != nil {
			t.Errorf("%d: %s", i, err)
		}
		if err := userPatchValid(&up); err != nil {
			if err.Error() != tt.expectedStatus {
				t.Errorf("%d: expected %s but got=%s", i, tt.expectedStatus, err)
			}
		}

	}
}
