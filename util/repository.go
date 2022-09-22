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
	if err != nil {
		log.Fatalln(err)
	}
	if db == nil {
		log.Fatalln("db nil")
	}
	err = createStruct(db)
	if err != nil {
		log.Fatalln(err)
	}
	return &Repository{db}
}

func (r *Repository) Close() {
	r.Db.Close()
}

func (r *Repository) InsertOrReplaceCards(items []models.Card) error {
	sql := `
	INSERT OR REPLACE INTO Cards(
		Code,
		Since,
		Until,
		Enabled
	) values(?, ?, ?, ?)
	`
	stmt, err := r.Db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		_, err2 := stmt.Exec(item.Code, item.Since, item.Until, item.Enabled)
		if err2 != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) ValidCard(code string) (bool, error) {
	sql := `
	SELECT COUNT(*)
	FROM Cards
	WHERE Code = ?
	AND CURRENT_TIMESTAMP BETWEEN Since AND Until
	AND Enabled = 1
	`
	stmt, err := r.Db.Prepare(sql)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	rows, err := r.Db.Query(sql, code)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	count := 0

	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return false, err
		}
	}
	return (count == 1), nil
}

func (r *Repository) InsertAccess(item *models.Access) error {
	sql := `
	INSERT INTO Access(
		Code1,
		Code2,
		AccessDate,
		Zone,
		Event
	) values(?, ?, ?, ?, ?)
	`
	stmt, err := r.Db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(item.Code1, item.Code2, item.AccessDate, item.Zone, item.Event)
	if err2 != nil {
		return err
	}
	return nil
}

func TestCards() []models.Card {
	return []models.Card{
		{Code: "001-42070AED-6E95-4586-91D0-F90BB10D1B7F", Since: time.Now(), Until: time.Now().Add(30 * time.Hour).UTC(), Enabled: true},
		{Code: "001-C343B902-AF36-4A5B-A3A8-C35A11C138E1", Since: time.Now(), Until: time.Now().Add(60 * time.Hour).UTC(), Enabled: false},
		{Code: "001-3F73A575-387C-45F5-9F48-396C5456C268", Since: time.Now(), Until: time.Now().Add(60 * time.Hour).UTC(), Enabled: true},
	}
}

func createStruct(Db *sql.DB) error {
	// create table if not exists
	sql_table := `
		CREATE TABLE IF NOT EXISTS Cards(
			Code VARCHAR(40) PRIMARY KEY,
			Since DATETIME,
			Until DATETIME,
			Enabled BOOL
		);`
	_, err := Db.Exec(sql_table)
	if err != nil {
		return err
	}

	sql_table = `
		CREATE TABLE IF NOT EXISTS Access(
			Code1 VARCHAR(40),
			Code2 VARCHAR(40),
			AccessDate DATETIME,
			Zone VARCHAR(2),
			Event VARCHAR(1)
		);`
	_, err = Db.Exec(sql_table)
	if err != nil {
		return err
	}
	return nil
}
