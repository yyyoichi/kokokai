package main

import (
	"encoding/json"
	"fmt"
	"kokokai/server/db"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/appengine/v2"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	loadEnv()
	http.HandleFunc("/", handle)
	http.HandleFunc("/daykyoki", handleDayKyoki)
	appengine.Main()
}
func handle(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("ENV")
	json.NewEncoder(w).Encode(Response{Status: "ok", Message: "Hello world! " + env})
}
func handleDayKyoki(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	dateString := r.URL.Query().Get("d")
	fmt.Println(dateString)
	regex := regexp.MustCompile(`[0-9]{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])`)
	if !regex.MatchString(dateString) {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(Response{Status: "Bad Request", Message: "No date"})
		return
	}
	connection, err := db.GetDatabase()
	if err != nil {
		w.WriteHeader(501)
		json.NewEncoder(w).Encode(Response{Status: "Service Unavailable", Message: "unconnect db"})
		return
	}
	defer connection.Close()
	kyoki := db.New(dateString, connection)
	println(kyoki.Kyoki)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(kyoki)
}

func loadEnv() {
	if !appengine.IsAppEngine() {
		currentDir, _ := os.Getwd()
		envPath := strings.ReplaceAll(filepath.Join(currentDir, ".env"), "\\", "/")
		err := godotenv.Load(envPath)
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}
}
