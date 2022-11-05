package db

import (
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
	db := GetDatabase()
	defer db.Close()
	if db == nil {
		t.Error("cannot connect database")
	}
}

func TestGetDatabase(t *testing.T) {
	db := GetDatabase()
	selectStmt := `select * from kyokiday limit 5`
	rows, er := db.Query(selectStmt)
	defer db.Close()
	if er != nil {
		t.Error("error")
	}
	for rows.Next() {
		var pk int64
		err := rows.Scan(&pk)
		if err != nil {
			t.Error("error scan")
		}
		t.Logf("pk:%d", pk)
	}
}
