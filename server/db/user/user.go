package user

import (
	"bytes"
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

type ColumnName string

const (
	PK        ColumnName = "pk"
	ID        ColumnName = "id"
	NAME      ColumnName = "name"
	EMAIL     ColumnName = "email"
	PASS      ColumnName = "pass"
	LOGIN_AT  ColumnName = "login_at"
	UPDATE_AT ColumnName = "update_at"
	CREATE_AT ColumnName = "create_at"
)

func (u *User) getField(cn ColumnName) interface{} {
	switch cn {
	case PK:
		return u.Pk
	case ID:
		return u.Id
	case NAME:
		return u.Name
	case EMAIL:
		return u.Email
	case PASS:
		return u.Pass
	case LOGIN_AT:
		return u.LoginAt
	case UPDATE_AT:
		return u.UpdateAt
	case CREATE_AT:
		return u.CreateAt
	default:
		return nil
	}
}

func (u *User) Create() error {
	if u.Pass == "" {
		return fmt.Errorf("empty pass")
	}
	if u.Id == "" {
		u.Id = newId()
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
		return fmt.Errorf("%s is exists", u.Id)
	}
	s := `INSERT INTO usr (id, name, pass) VALUES($1, $2, $3) RETURNING pk`
	var pk int64
	err = conn.QueryRow(s, u.Id, u.Id, u.Pass).Scan(&pk)
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

// IDから取得
func (u *User) GetById() error {
	if u.Id == "" {
		return fmt.Errorf("empty id")
	}
	conn, err := db.GetDatabase()
	if err != nil {
		return err
	}
	defer conn.Close()
	// この辺どうにかしたい
	var pk sql.NullInt64
	var id sql.NullString
	var name sql.NullString
	var email sql.NullString
	var pass sql.NullString
	var loginAt sql.NullTime
	var updateAt sql.NullTime
	var createAt sql.NullTime
	s := `SELECT * FROM usr WHERE id=$1`
	err = conn.QueryRow(s, u.Id).Scan(&pk, &id, &name, &email, &pass, &loginAt, &updateAt, &createAt)
	if err != nil {
		return err
	}
	u.Pk, u.Id, u.Name, u.Email, u.Pass = db.N2i(pk), db.N2s(id), db.N2s(name), db.N2s(email), db.N2s(pass)
	u.LoginAt, u.UpdateAt, u.CreateAt = db.N2t(loginAt), db.N2t(updateAt), db.N2t(createAt)
	if u.loginstamp(conn) != nil {
		return err
	}
	return nil
}

// パスワード確認後取得
func (u *User) GetByPass() error {
	if u.Pass == "" {
		return fmt.Errorf("empty password")
	}
	if u.Id == "" {
		return fmt.Errorf("empty id")
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
	s := `SELECT * FROM usr WHERE id=$1 and pass=$2`
	err = conn.QueryRow(s, u.Id, u.Pass).Scan(&pk, &id, &name, &email, &pass, &loginAt, &updateAt, &createAt)
	if err != nil {
		exists, err := u.exists(conn)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("wrong-pass")
		} else {
			return fmt.Errorf("no-id")
		}
	}
	u.Pk, u.Id, u.Name, u.Email, u.Pass = db.N2i(pk), db.N2s(id), db.N2s(name), db.N2s(email), db.N2s(pass)
	u.LoginAt, u.UpdateAt, u.CreateAt = db.N2t(loginAt), db.N2t(updateAt), db.N2t(createAt)
	if u.loginstamp(conn) != nil {
		return err
	}
	return nil
}

// 認証済み
func (u *User) Update(columns []ColumnName) error {
	if len(columns) == 0 {
		return nil
	}
	if u.Id == "" {
		return fmt.Errorf("empty id")
	}
	conn, err := db.GetDatabase()
	if err != nil {
		return err
	}
	defer conn.Close()

	// プレスホルダー
	var ph []interface{}
	// sql文
	var s bytes.Buffer
	now := time.Now()
	ph = append(ph, now)
	s.WriteString("UPDATE usr SET update_at=$1")
	// アップデート対象
	for _, cn := range columns {
		ph = append(ph, u.getField(cn))
		s.WriteString(fmt.Sprintf(", %s=$%d", cn, len(ph)))
	}
	ph = append(ph, u.Id)
	s.WriteString(fmt.Sprintf(" WHERE id=$%d", len(ph)))
	ss := s.String()
	_, err = conn.Exec(ss, ph...)
	if err != nil {
		return err
	}
	u.UpdateAt = now
	return nil
}

func (u *User) Delete() error {
	if u.Id == "" {
		return fmt.Errorf("empty id")
	}
	conn, err := db.GetDatabase()
	if err != nil {
		return err
	}
	defer conn.Close()
	s := `DELETE FROM usr WHERE id=$1`
	_, err = conn.Exec(s, u.Id)
	if err != nil {
		return err
	}
	u = &User{Id: ""}
	return nil
}

func (u *User) exists(conn *sql.DB) (bool, error) {
	if u.Id == "" {
		return false, fmt.Errorf("empty id")
	}
	s := `SELECT id FROM usr WHERE id=$1`
	rows, err := conn.Query(s, u.Id)
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
