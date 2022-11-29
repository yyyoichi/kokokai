package handle

import (
	"encoding/json"
	"fmt"
	"kokokai/server/db"
	"net/http"
	"regexp"
)

func DayKyoki(w http.ResponseWriter, r *http.Request) {
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
