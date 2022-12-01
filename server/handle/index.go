package handle

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
}

func (res *Response) resError(w *http.ResponseWriter) {
	json, err := json.Marshal(res)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(*w, string(json), http.StatusBadRequest)
}

func Index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Response{Status: "ok"})
}

func customAlphanumary(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	rg := regexp.MustCompile(`([0-9].*[a-zA-Z]|[a-zA-Z].*[0-9])`)
	return rg.MatchString(s)
}
