package util

import (
	"database/sql"
	"log"

	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sehogas/qr-reader/models"
)

type Repository struct {
	Db     *sql.DB
	Config *Config
}

type Config struct {
	LastUpdateCards time.Time
	Anulados        bool
}

func NewRepository(driverName string, filepath string) *Repository {

	defaultConfig := &Config{
		LastUpdateCards: time.Now().AddDate(-50, 0, 0),
		Anulados:        false,
	}

	db, err := sql.Open(driverName, filepath)
	if err != nil {
		log.Fatal(err)
	}
	if db == nil {
		log.Fatal("db nil")
	}

	err = createStruct(db)
	if err != nil {
		log.Fatalln(err)
	}

	config, err := readConfig(db)
	if err != nil {
		log.Fatalln(err)
	}

	if (config == Config{}) {
		config = *defaultConfig
		err = insertDefaultConfig(db, defaultConfig)
		if err != nil {
			log.Fatal("Set default config: ", err)
		}
	}

	return &Repository{Db: db, Config: &config}
}

func (r *Repository) Close() {
	r.Db.Close()
}

func createStruct(Db *sql.DB) error {
	sql_table := `
		CREATE TABLE IF NOT EXISTS Config(
			LastUpdateCards DATETIME
		);`
	_, err := Db.Exec(sql_table)
	if err != nil {
		return err
	}

	sql_table = `
		CREATE TABLE IF NOT EXISTS Cards(
			Code VARCHAR(40) PRIMARY KEY,
			DateFrom DATETIME,
			DateTo DATETIME,
			Enabled BOOL,
			Photo VARCHAR(255),
			Deleted BOOL
		);`
	_, err = Db.Exec(sql_table)
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

func insertDefaultConfig(Db *sql.DB, config *Config) error {
	sql := `
	INSERT INTO Config(
		LastUpdateCards
	) values(?)
	`
	stmt, err := Db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(config.LastUpdateCards)
	if err != nil {
		return err
	}
	return nil
}

func readConfig(Db *sql.DB) (Config, error) {
	var config Config
	sql := `
	SELECT LastUpdateCards
	FROM Config
	`
	stmt, err := Db.Prepare(sql)
	if err != nil {
		return config, err
	}
	defer stmt.Close()

	rows, err := Db.Query(sql)
	if err != nil {
		return config, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&config.LastUpdateCards)
		if err != nil {
			return config, err
		}
	}
	return config, nil
}

func (r *Repository) SyncCards(items []models.Card) error {
	sql := `
		INSERT OR REPLACE INTO Cards(
			Code,
			DateFrom,
			DateTo,
			Enabled,
			Photo,
			Deleted
		) VALUES (?, ?, ?, ?, ?, ?)
	`
	stmt, err := r.Db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, i := range items {
		_, err = stmt.Exec(i.Code, i.DateFrom, i.DateTo, i.Enabled, i.Photo, i.Deleted)
		if err != nil {
			return err
		}
	}

	sql = `
		DELETE FROM Cards WHERE Deleted = 1
		`
	stmt2, err2 := r.Db.Prepare(sql)
	if err2 != nil {
		return err2
	}
	defer stmt2.Close()

	_, err = stmt2.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) TotalCards() (int64, error) {
	var count int64 = 0
	sql := `
	SELECT COUNT(*)
	FROM Cards
	`
	stmt, err := r.Db.Prepare(sql)
	if err != nil {
		return count, err
	}
	defer stmt.Close()

	row := r.Db.QueryRow(sql)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) ValidCard(code string) (bool, error) {
	sql := `
	SELECT COUNT(*)
	FROM Cards
	WHERE Code = ?
	AND CURRENT_TIMESTAMP BETWEEN DateFrom AND DateTo
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

	_, err = stmt.Exec(item.Code1, item.Code2, item.AccessDate, item.Zone, item.Event)
	if err != nil {
		return err
	}
	return nil
}

func TestCards() []models.Card {
	return []models.Card{
		{Code: "001-42070AED-6E95-4586-91D0-F90BB10D1B7F", DateFrom: time.Now(), DateTo: time.Now().Add(30 * time.Hour).UTC(), Enabled: true, Photo: "", Deleted: false},
		{Code: "001-C343B902-AF36-4A5B-A3A8-C35A11C138E1", DateFrom: time.Now(), DateTo: time.Now().Add(60 * time.Hour).UTC(), Enabled: false, Photo: "", Deleted: false},
		{Code: "001-3F73A575-387C-45F5-9F48-396C5456C268", DateFrom: time.Now(), DateTo: time.Now().Add(60 * time.Hour).UTC(), Enabled: true, Photo: "", Deleted: false},
		{Code: "002-94E5801E-69BE-4268-9B18-F5C4CB7C5187", DateFrom: time.Now(), DateTo: time.Now().Add(60 * time.Hour).UTC(), Enabled: true, Photo: "", Deleted: false},
		{Code: "002-94E5801E-69BE-4268-9B18-F5C4CB7C5181", DateFrom: time.Now(), DateTo: time.Now().Add(60 * time.Hour).UTC(), Enabled: false, Photo: "", Deleted: false},
	}
}
