package handle

import (
	"encoding/json"
	"net/http"
	"os"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("ENV")
	json.NewEncoder(w).Encode(Response{Status: "ok", Message: "Hello world! " + env})
}
