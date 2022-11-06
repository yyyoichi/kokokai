package db

import (
	"database/sql"
	"strings"
	"testing"
)

func TestGetPsqlConn(t *testing.T) {
	conn := getPsqlConn()
	configs := strings.Split(conn, " ")
	for _, cnf := range configs {
		env := strings.Split(cnf, "=")
		key := env[0]
		val := env[1]
		if val == "" {
			t.Errorf("conn %s value==''", key)
		}
		t.Log(cnf)
	}
}
func TestDatabaseConnect(t *testing.T) {
	db, err := GetDatabase()
	if err != nil {
		t.Error("cannot connect database")
	}
	err = db.Ping()
	if err != nil {
		t.Error("failure")
	} else {
		t.Log("success")
	}
	defer db.Close()
}

func TestQueryDatabase(t *testing.T) {
	db, err := GetDatabase()
	if err != nil {
		t.Errorf("cannot connect db.")
	}
	selectStmt := `select * from kyokiday limit 5`
	rows, err := db.Query(selectStmt)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var pk sql.NullInt64
		var date sql.NullTime
		var createAt sql.NullTime
		err := rows.Scan(&pk, &date, &createAt)
		if err != nil {
			t.Error("error scan")
		}
		t.Log(pk)
		t.Log(date)
		t.Log(createAt)
	}
	defer db.Close()
}
