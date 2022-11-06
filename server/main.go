package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"google.golang.org/appengine/v2"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/", handle)
	appengine.Main()
	env := os.Getenv("ENV")
	fmt.Printf("env: %s", env)
}
func handle(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("ENV")
	json.NewEncoder(w).Encode(Response{Status: "ok", Message: "Hello world! " + env})
}
