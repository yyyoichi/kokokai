package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"google.golang.org/appengine/v2"
)

func getPsqlConn() string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	conf := fmt.Sprintf(
		"user=%s password=%s database=%s host=%s sslmode=disable",
		user, password, dbname, host,
	)
	return conf
}
func GetDatabase() (*sql.DB, error) {
	var db *sql.DB
	var err error
	if appengine.IsAppEngine() {
		db, err = sql.Open("pgx", getPsqlConn())
	} else {
		db, err = sql.Open("postgres", getPsqlConn())
	}
	if err != nil {
		return nil, err
	}
	return db, err
}
