package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

func getPsqlConn() string {
	currentDir, _ := os.Getwd()
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	sslrootcert := strings.ReplaceAll(filepath.Join(currentDir, os.Getenv("SSL_ROOT_CERT")), "\\", "/")
	sslkey := strings.ReplaceAll(filepath.Join(currentDir, os.Getenv("SSL_KEY")), "\\", "/")
	sslcert := strings.ReplaceAll(filepath.Join(currentDir, os.Getenv("SSL_CERT")), "\\", "/")
	psqlconn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=verify-ca sslrootcert=%s sslkey=%s sslcert=%s",
		host, port, user, password, dbname, sslrootcert, sslkey, sslcert,
	)
	return psqlconn
}
func GetDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", getPsqlConn())
	if err != nil {
		return nil, err
	}
	return db, err
}
