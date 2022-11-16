package db

import (
	"database/sql"
	"time"
)

func N2(v interface{}) interface{} {
	switch v := v.(type) {
	case sql.NullString:
		return n2s(v)
	case sql.NullInt64:
		return n2i(v)
	case sql.NullTime:
		return n2t(v)
	default:
		return nil
	}
}

func n2s(v sql.NullString) string {
	var s string
	if v.Valid {
		s = v.String
	}
	return s
}

func n2i(v sql.NullInt64) int64 {
	var i int64
	if v.Valid {
		i = v.Int64
	}
	return i
}

func n2t(v sql.NullTime) time.Time {
	var t time.Time
	if v.Valid {
		t = v.Time
	}
	return t
}
