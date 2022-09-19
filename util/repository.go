package util

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sehogas/qr-reader/models"
)

type Repository struct {
	Db *sql.DB
}

func NewRepository(driverName string, filepath string) *Repository {
	db, err := sql.Open(driverName, filepath)
	checkErr(err)
	if db == nil {
		panic("db nil")
	}
	createStruct(db)
	return &Repository{db}
}

func (r *Repository) Close() {
	r.Db.Close()
}

func (r *Repository) InsertOrReplaceCards(items []models.Card) {
	sql := `
	INSERT OR REPLACE INTO Cards(
		Code,
		Since,
		Until,
		Enabled
	) values(?, ?, ?, ?)
	`
	stmt, err := r.Db.Prepare(sql)
	checkErr(err)
	defer stmt.Close()

	for _, item := range items {
		_, err2 := stmt.Exec(item.Code, item.Since, item.Until, item.Enabled)
		checkErr(err2)
	}
}

func (r *Repository) ValidCard(code string) bool {
	sql := `
	SELECT COUNT(*)
	FROM Cards
	WHERE Code = ?
	AND CURRENT_TIMESTAMP BETWEEN Since AND Until
	AND Enabled = 1
	`
	stmt, err := r.Db.Prepare(sql)
	checkErr(err)
	defer stmt.Close()

	rows, err := r.Db.Query(sql, code)
	checkErr(err)
	defer rows.Close()

	count := 0

	if rows.Next() {
		err = rows.Scan(&count)
		checkErr(err)
	}
	return (count == 1)
}

func (r *Repository) InsertAccess(item models.Access) {
	sql := `
	INSERT INTO Access(
		Code,
		AccessDate,
		Zone,
		Event,
		Synchronized
	) values(?, ?, ?, ?, ?)
	`
	stmt, err := r.Db.Prepare(sql)
	checkErr(err)
	defer stmt.Close()

	_, err2 := stmt.Exec(item.Code, item.AccessDate, item.Zone, item.Event, item.Synchronized)
	checkErr(err2)
}

func TestCards() []models.Card {
	return []models.Card{
		{Code: "001-42070AED-6E95-4586-91D0-F90BB10D1B7F", Since: time.Now(), Until: time.Now().Add(30 * time.Hour).UTC(), Enabled: true},
		{Code: "001-C343B902-AF36-4A5B-A3A8-C35A11C138E2", Since: time.Now(), Until: time.Now().Add(60 * time.Hour).UTC(), Enabled: true},
	}
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func createStruct(Db *sql.DB) {
	// create table if not exists
	sql_table := `
		CREATE TABLE IF NOT EXISTS Cards(
			Code VARCHAR(40) PRIMARY KEY,
			Since DATETIME,
			Until DATETIME,
			Enabled BOOL
		);`
	_, err := Db.Exec(sql_table)
	checkErr(err)

	sql_table = `
		CREATE TABLE IF NOT EXISTS Access(
			Code VARCHAR(40),
			AccessDate DATETIME,
			Zone VARCHAR(2),
			Event VARCHAR(1),
			Synchronized BOOL
		);`
	_, err = Db.Exec(sql_table)
	checkErr(err)
}
