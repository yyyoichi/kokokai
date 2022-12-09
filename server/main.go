package main

import (
	"kokokai/server/handle"
	"kokokai/server/handle/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/csrf"
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

	s := r.PathPrefix("/sessions/").Subrouter()
	s.Use(middleware.MiddlewareAuth)
	s.HandleFunc("/", handle.UserSessionFunc).Methods("GET")

	u := r.PathPrefix("/users/").Subrouter()
	u.Use(middleware.MiddlewareAuth)
	csrfMiddleware := getCSRFMiddleware()
	u.Use(csrfMiddleware)
	u.HandleFunc("/{userId}", handle.UserFunc).Methods("PATCH")

	http.Handle("/", r)
}

func getCSRFMiddleware() mux.MiddlewareFunc {
	ev := os.Getenv("DEV")
	s := []byte(os.Getenv("CSRF_SECRET"))
	var md mux.MiddlewareFunc
	switch ev {
	case "DEV":
		md = csrf.Protect(s)
	case "STG":
		md = csrf.Protect(s, csrf.TrustedOrigins([]string{"collokaistg.yyyoichi.com"}))
	case "PRO":
		md = csrf.Protect(s, csrf.TrustedOrigins([]string{"collokai.yyyoichi.com"}))
	}
	return md
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
