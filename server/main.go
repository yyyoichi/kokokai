package main

import (
	"kokokai/server/handle"
	"kokokai/server/handle/middleware"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/appengine/v2"
)

func main() {
	loadEnv()
	handler()
	appengine.Main()
}

func handler() {
	r := mux.NewRouter()
	r.HandleFunc("/", handle.Index).Methods("GET")
	r.HandleFunc("/daykyoki", handle.DayKyoki).Methods("GET")
	r.HandleFunc("/login", handle.LoginFunc).Methods("POST")
	r.HandleFunc("/signup", handle.SignUpFunc).Methods("POST")
	u := r.PathPrefix("/users/").Subrouter()
	u.Use(middleware.MiddlewareAuth)
	u.HandleFunc("/users/{userId}", handle.UserFunc).Methods("PATCH")
}

func loadEnv() {
	if appengine.IsAppEngine() {
		return
	}
	currentDir, _ := os.Getwd()
	envPath := strings.ReplaceAll(filepath.Join(currentDir, ".env"), "\\", "/")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
