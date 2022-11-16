package auth

import (
	"database/sql"
	"fmt"
	"kokokai/server/db"
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

func (u *User) CreateUser(conn *sql.DB) (*User, error) {
	s := `INSERT INTO usr (id, name, email, pass) VALUES($1, $2, $3, $4)`
	res, err := conn.Exec(s, u.Id, u.Name, u.Email, u.Pass)
	if err != nil {
		return nil, err
	}
	pk, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.Pk = pk
	return u, nil
}

func (u *User) GetUser(conn *sql.DB) (*User, error) {
	s := `SELECT * FROM usr WHERE email=$1, pass=$1`
	rows, err := conn.Query(s, u.Email, u.Pass)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		var pk sql.NullInt64
		var id sql.NullString
		var name sql.NullString
		var email sql.NullString
		var pass sql.NullString
		var loginAt sql.NullTime
		var updateAt sql.NullTime
		var createAt sql.NullTime
		rows.Scan(&pk, &id, &name, &email, &pass, &loginAt, &updateAt, &createAt)
		u.Pk, u.Id, u.Name, u.Email, u.Pass, u.LoginAt, u.UpdateAt, u.CreateAt = db.N2i(pk), db.N2s(id), db.N2s(name), db.N2s(email), db.N2s(pass), db.N2t(loginAt), db.N2t(updateAt), db.N2t(createAt)
		err := u.loginstamp(conn)
		if err != nil {
			return u, err
		}
		return u, nil
	}
	return nil, fmt.Errorf("no-data")
}

func (u *User) loginstamp(conn *sql.DB) error {
	s := `UPDATE usr WHERE pk=$1 SET login_at=$2`
	now := time.Now()
	_, err := conn.Exec(s, u.Pk, now)
	if err != nil {
		return err
	}
	u.LoginAt = now
	return nil
}