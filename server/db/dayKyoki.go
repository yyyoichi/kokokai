package db

import (
	"database/sql"
)

type Kyoki struct {
	Words []string `json:"words"`
	Pk    int64    `json:"pk"`
	Freq  int64    `json:"freq"`
}

type DayKyoki struct {
	db         *sql.DB  `json:"-"`
	DateString string   `json:"date"`
	Kyoki      []*Kyoki `json:"kyoki"`
}

func New(dateString string, db *sql.DB) *DayKyoki {
	return &DayKyoki{db: db, DateString: dateString, Kyoki: make([]*Kyoki, 0)}
}

func (d *DayKyoki) Get() *DayKyoki {
	d.Kyoki = d.getKyoki()
	for _, k := range d.Kyoki {
		kyokiPk := k.Pk
		k.Words = d.GetKyokiItem(kyokiPk)
	}
	return d
}
func (d *DayKyoki) getKyoki() []*Kyoki {
	selectStmt := `
	SELECT k.pk AS pk, k.freq AS freq
	FROM kyokiday kd
	JOIN kyoki k ON kd.pk = k.kyokiday
	WHERE kd.date = $1
	ORDER BY k.freq DESC
	LIMIT 30
	`
	rows, err := d.db.Query(selectStmt, d.DateString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var kyokiList []*Kyoki
	for rows.Next() {
		var pk sql.NullInt64
		var freq sql.NullInt64
		err := rows.Scan(&pk, &freq)
		if err != nil {
			panic(err)
		}
		kyokiList = append(kyokiList, &Kyoki{Pk: n2i(pk), Freq: n2i(freq), Words: make([]string, 0)})
	}
	return kyokiList
}
func (d *DayKyoki) GetKyokiItem(kyokiPk int64) []string {
	selectStmt := `
	SELECT w.word AS word
	FROM kyokiitem ki JOIN word w ON ki.word = w.code
	WHERE ki.kyoki = $1
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