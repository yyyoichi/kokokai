package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/appengine/v2"
)

func getPsqlConn() string {
	if !appengine.IsAppEngine() {
		err := godotenv.Load("config/dev/.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	currentDir, _ := os.Getwd()
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	sslrootcert := filepath.Join(currentDir, os.Getenv("SSL_ROOT_CERT"))
	sslkey := filepath.Join(currentDir, os.Getenv("SSL_CERT"))
	sslcert := filepath.Join(currentDir, os.Getenv("SSL_KEY"))
	psqlconn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=verify-full sslrootcert=%s sslkey=%s sslcert=%s",
		host, port, user, password, dbname, sslrootcert, sslkey, sslcert,
	)
	return psqlconn
}
func GetDatabase() *sql.DB {

	db, err := sql.Open("postgres", getPsqlConn())
	CheckError(err)
	return db
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
