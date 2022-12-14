package db

import (
	"database/sql"
	"time"
)

func N2s(v sql.NullString) string {
	var s string
	if v.Valid {
		s = v.String
	}
	return s
}

func N2i(v sql.NullInt64) int64 {
	var i int64
	if v.Valid {
		i = v.Int64
	}
	return i
}

func N2t(v sql.NullTime) time.Time {
	var t time.Time
	if v.Valid {
		t = v.Time
	}
	return t
}
