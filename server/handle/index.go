package handle

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status string `json:"status"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Response{Status: "ok"})
}
