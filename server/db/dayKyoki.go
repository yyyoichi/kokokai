package db

import (
	"database/sql"
	"strings"
)

type Kyoki struct {
	Words []string `json:"words"`
	Pk    int64    `json:"pk"`
	Freq  int64    `json:"freq"`
}

type DayKyoki struct {
	DateString string   `json:"date"`
	Kyoki      []*Kyoki `json:"kyoki"`
}

func New(dateString string, db *sql.DB) *DayKyoki {
	selectStmt := `
	SELECT 
	k.pk AS pk, 
	k.freq AS freq,
	array_to_string(
		ARRAY(
			SELECT w.word
			FROM word w
			JOIN kyokiitem ki ON ki.kyoki = k.pk
			WHERE w.code = ki.word
		), ','
	) AS words
	FROM kyokiday kd
	JOIN kyoki k ON kd.pk = k.kyokiday
	WHERE kd.date = $1
	ORDER BY k.freq DESC
	LIMIT 30
	`
	rows, err := db.Query(selectStmt, dateString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var kyoki []*Kyoki
	for rows.Next() {
		var pk sql.NullInt64
		var freq sql.NullInt64
		var words sql.NullString
		rows.Scan(&pk, &freq, &words)
		k := &Kyoki{Pk: n2i(pk), Freq: n2i(freq), Words: strings.Split(n2s(words), ",")}
		kyoki = append(kyoki, k)
	}
	defer rows.Close()
	return &DayKyoki{DateString: dateString, Kyoki: kyoki}
}
