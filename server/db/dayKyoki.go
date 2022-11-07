package db

import (
	"database/sql"
)

type Kyoki struct {
	words []string `json: "words"`
	pk    int64    `json: "pk"`
	freq  int64    `json: "freq"`
}

type DayKyoki struct {
	db         *sql.DB `json: "-"`
	dateString string  `json: "date"`
	kyoki      []Kyoki `json: "kyoki"`
}

func New(dateString string, db *sql.DB) DayKyoki {
	return DayKyoki{db: db, dateString: dateString, kyoki: make([]Kyoki, 0)}
}

func (d *DayKyoki) Get() *DayKyoki {
	d.kyoki = d.getKyoki()
	return d
}
func (d *DayKyoki) getKyoki() []Kyoki {
	selectStmt := `
	SELECT k.pk AS pk, k.freq AS freq
	FROM kyokiday kd
	JOIN kyoki k ON kd.pk = kyoki.kyokiday
	WHERE kd.pk ='?'
	ORDER BY k.req DESC
	LIMIT 30
	`
	rows, err := d.db.Query(selectStmt, d.dateString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var kyokiList []Kyoki
	for rows.Next() {
		var pk sql.NullInt64
		var freq sql.NullInt64
		err := rows.Scan(&pk, &freq)
		if err != nil {
			panic(err)
		}
		kyokiPk := n2i(pk)
		words := d.GetKyokiItem(kyokiPk)
		kyokiList = append(kyokiList, Kyoki{pk: kyokiPk, freq: n2i(freq), words: words})
	}
	return kyokiList
}
func (d *DayKyoki) GetKyokiItem(kyokiPk int64) []string {
	selectStmt := `
	SELECT w.word AS word
	FORM kyokiitem ki JOIN word w ON ki.word = w.code
	WHERE ki.kyoki = ?
	`
	rows, err := d.db.Query(selectStmt, kyokiPk)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var words []string
	for rows.Next() {
		var word sql.NullString
		err := rows.Scan(&word)
		if err != nil {
			panic(err)
		}
		words = append(words, n2s(word))
	}
	return words
}
