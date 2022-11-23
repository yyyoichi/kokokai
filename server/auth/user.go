package auth

import (
	"database/sql"
	"fmt"
	"kokokai/server/db"
	"math/rand"
	"strings"
	"time"
)

type User struct {
	Pk       int64
	Id       string
	Name     string
	Email    string
	Pass     string
	LoginAt  time.Time
	UpdateAt time.Time
	CreateAt time.Time
}

func (u *User) Create() error {
	if u.Pass == "" || u.Email == "" {
		return fmt.Errorf("empty pass or email")
	}
	conn, err := db.GetDatabase()
	if err != nil {
		return err
	}
	defer conn.Close()
	exists, err := u.exists(conn)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%s is exists", u.Email)
	}
	s := `INSERT INTO usr (id, name, email, pass) VALUES($1, $2, $3, $4) RETURNING pk`
	u.Id = newId()
	var pk int64
	err = conn.QueryRow(s, u.Id, u.Id, u.Email, u.Pass).Scan(&pk)
	if err != nil {
		return err
	}
	u.Pk = pk
	return nil
}

func newId() string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var b strings.Builder
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20; i++ {
		b.WriteByte(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func (u *User) Get() error {
	if u.Pass == "" {
		return fmt.Errorf("empty password")
	}
	if u.Email == "" && u.Id == "" {
		return fmt.Errorf("empty email or id")
	}
	conn, err := db.GetDatabase()
	if err != nil {
		return err
	}
	defer conn.Close()
	var pk sql.NullInt64
	var id sql.NullString
	var name sql.NullString
	var email sql.NullString
	var pass sql.NullString
	var loginAt sql.NullTime
	var updateAt sql.NullTime
	var createAt sql.NullTime
	s := `SELECT * FROM usr WHERE (email=$1 and pass=$2) or (id=$3 and pass=$2)`
	err = conn.QueryRow(s, u.Email, u.Pass, u.Id).Scan(&pk, &id, &name, &email, &pass, &loginAt, &updateAt, &createAt)
	if err != nil {
		exists, err := u.exists(conn)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("wrong-pass")
		} else {
			return fmt.Errorf("no-email")
		}
	}
	u.Pk, u.Id, u.Name, u.Email, u.Pass = db.N2i(pk), db.N2s(id), db.N2s(name), db.N2s(email), db.N2s(pass)
	u.LoginAt, u.UpdateAt, u.CreateAt = db.N2t(loginAt), db.N2t(updateAt), db.N2t(createAt)
	if u.loginstamp(conn) != nil {
		return err
	}
	return nil
}

func (u *User) Delete() error {
	if u.Pk == 0 || u.Email == "" {
		return fmt.Errorf("no pk or email")
	}
	conn, err := db.GetDatabase()
	if err != nil {
		return err
	}
	defer conn.Close()
	s := `DELETE FROM usr WHERE pk=$1 or email=$2`
	_, err = conn.Exec(s, u.Pk, u.Email)
	if err != nil {
		return err
	}
	u = &User{Pk: 0, Email: "", Pass: ""}
	return nil
}

func (u *User) exists(conn *sql.DB) (bool, error) {
	if u.Email == "" {
		return false, fmt.Errorf("empty email")
	}
	s := `SELECT pk FROM usr WHERE email=$1`
	rows, err := conn.Query(s, u.Email)
	if err != nil {
		return false, err
	}
	return rows.Next(), nil
}

func (u *User) loginstamp(conn *sql.DB) error {
	if u.Pk == 0 {
		return fmt.Errorf("empty pk")
	}
	s := `UPDATE usr WHERE pk=$1 SET login_at=$2`
	now := time.Now()
	_, err := conn.Exec(s, u.Pk, now)
	if err != nil {
		return err
	}
	u.LoginAt = now
	return nil
}
