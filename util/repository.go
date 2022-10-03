package util

import (
	"database/sql"
	"fmt"
	"log"

	//	"net/url"
	"time"

	//_ "github.com/mattn/go-sqlite3"

	_ "github.com/xeodou/go-sqlcipher"

	"github.com/sehogas/qr-reader/models"
)

type Repository struct {
	Db *sql.DB
}

func NewRepository(driverName string, filepath string) *Repository {

	//	password := "lalala"

	//dbnameWithDSN := filepath + fmt.Sprintf("?_pragma_key=%s&_pragma_cipher_page_size=4096", url.QueryEscape(password))
	//db, err := sql.Open("sqlite3", dbnameWithDSN)

	//db, err := sql.Open(driverName, filepath)

	/*
		sql.Register("sqlite3_log", &sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				log.Printf("Auth enabled: %v\n", conn.AuthEnabled())
				return nil
			},
		})

		// This is usual DB stuff (except with our sqlite3_log driver)
		db, err := sql.Open("sqlite3_log", fmt.Sprintf("file:%s?_auth&_auth_user=admin&_auth_pass=lalala&_auth_crypt=SSHA512&_auth_salt=233446", filepath))
	*/

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_key=123456", filepath))
	if err != nil {
		log.Fatal(err)
	}
	if db == nil {
		log.Fatalln("db nil")
	}

	p := "PRAGMA key = '123456';"
	_, err = db.Exec(p)
	if err != nil {
		log.Fatal(err)
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
		{Code: "002-94E5801E-69BE-4268-9B18-F5C4CB7C5187", Since: time.Now(), Until: time.Now().Add(60 * time.Hour).UTC(), Enabled: true},
		{Code: "002-94E5801E-69BE-4268-9B18-F5C4CB7C5181", Since: time.Now(), Until: time.Now().Add(60 * time.Hour).UTC(), Enabled: false},
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
