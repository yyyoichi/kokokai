package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func testLoadEnv() {
	currentDir, _ := os.Getwd()
	envPath := strings.ReplaceAll(filepath.Join(currentDir, "config/dev/test_db.env"), "\\", "/")
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func TestGetPsqlConn(t *testing.T) {
	testLoadEnv()
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
	testLoadEnv()
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
	testLoadEnv()
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

func TestNullXXConvert(t *testing.T) {
	testLoadEnv()
	db, err := GetDatabase()
	if err != nil {
		t.Errorf("cannot connect db.")
	}
	rows, err := db.Query(`select * from kyokiday where pk > 0 limit 1`)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var pk sql.NullInt64
		var date sql.NullTime
		var createAt sql.NullTime
		err := rows.Scan(&pk, &date, &createAt)
		if err != nil {
			t.Error("error scan")
		}
		testingInteger(t, pk)
		testingTime(t, date)
		testingTime(t, createAt)
	}
	rows, err = db.Query(`select * from word where word is not null limit 1`)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var code sql.NullInt64
		var word sql.NullString
		err := rows.Scan(&code, &word)
		if err != nil {
			t.Error("error scan")
		}
		testingInteger(t, code)
		testingString(t, word)
	}
	defer db.Close()
}

func TestDaykyoki(t *testing.T) {
	testLoadEnv()
	db, err := GetDatabase()
	if err != nil {
		t.Errorf("cannot connect db.")
	}
	defer db.Close()
	dateString := "2022-10-26"
	kyoki := New(dateString, db)
	kyokiList := kyoki.Kyoki
	if len(kyokiList) != 30 {
		t.Errorf("kyokiList len=%d", len(kyokiList))
	}
	t.Logf("pk: %d, freq: %d", kyokiList[0].Pk, kyokiList[0].Freq)
}

func testingString(t *testing.T, v interface{}) {
	nv, ok := v.(sql.NullString)
	if !ok {
		t.Errorf("exp not sql.NullString. got=%T", v)
	}
	s := n2s(nv)
	t.Log(s)
}
func testingInteger(t *testing.T, v interface{}) {
	nv, ok := v.(sql.NullInt64)
	if !ok {
		t.Errorf("exp not sql.NullInt64. got=%T", v)
	}
	s := n2i(nv)
	t.Log(s)
}
func testingTime(t *testing.T, v interface{}) {
	nv, ok := v.(sql.NullTime)
	if !ok {
		t.Errorf("exp not sql.NullTime. got=%T", v)
	}
	s := n2t(nv)
	t.Log(s)
}
