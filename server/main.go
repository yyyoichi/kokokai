package main

import (
	"kokokai/server/handle"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/appengine/v2"
)

func main() {
	loadEnv()
	handler()
	appengine.Main()
}

func handler() {
	http.HandleFunc("/", handle.Index)
	http.HandleFunc("/daykyoki", handle.DayKyoki)
	http.HandleFunc("/login", handle.LoginFunc)
	http.HandleFunc("/signup", handle.SignUpFunc)
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
