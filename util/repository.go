package util

import (
	"database/sql"
	"log"

	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sehogas/qr-reader/models"
)

type Repository struct {
	Db     *sql.DB
	Config *Config
}

type Config struct {
	LastUpdateCards time.Time
}

func NewRepository(driverName string, filepath string) *Repository {

	defaultConfig := &Config{
		LastUpdateCards: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
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
		err = insertOrReplaceDefaultConfig(db, defaultConfig)
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
			UUID VARCHAR(36),
			Code1 VARCHAR(40),
			Code2 VARCHAR(40),
			AccessDate DATETIME,
			Zone VARCHAR(2),
			Event VARCHAR(1),
			SyncDate DATETIME
		);`
	_, err = Db.Exec(sql_table)
	if err != nil {
		return err
	}

	return nil
}

func insertOrReplaceDefaultConfig(Db *sql.DB, config *Config) error {
	sql := `
	INSERT OR REPLACE INTO Config(
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

func (r *Repository) NewUUID() string {
	return uuid.New().String()
}

func (r *Repository) SyncCards(items []models.Card, serverTime time.Time) error {
	if len(items) > 0 {
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

		_, err2 = stmt2.Exec()
		if err2 != nil {
			return err2
		}

		sql = `
		UPDATE Config SET LastUpdateCards = ?
		`
		stmt3, err3 := r.Db.Prepare(sql)
		if err3 != nil {
			return err3
		}
		defer stmt3.Close()

		_, err3 = stmt3.Exec(serverTime)
		if err3 != nil {
			return err3
		}

		r.Config.LastUpdateCards = serverTime
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

func (r *Repository) TotalEarrings() (int64, error) {
	var count int64 = 0
	sql := `
	SELECT COUNT(*)
	FROM Access
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

func (r *Repository) TotalAccessToSync() (int64, error) {
	var count int64 = 0
	sql := "SELECT COUNT(*)	FROM Access	WHERE SyncDate IS NULL"

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

func (r *Repository) GetAccessToSync(date time.Time) ([]models.Access, error) {
	var items []models.Access
	sql := "UPDATE Access SET SyncDate = ? WHERE SyncDate IS NULL"
	stmt, err := r.Db.Prepare(sql)
	if err != nil {
		return items, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(date)
	if err != nil {
		return items, err
	}

	sql = "SELECT UUID, Code1, Code2, AccessDate, Zone, Event, SyncDate FROM Access WHERE SyncDate = ?"
	stmt2, err2 := r.Db.Prepare(sql)
	if err2 != nil {
		return items, err2
	}
	defer stmt2.Close()

	rows, err2 := r.Db.Query(sql, date)
	if err2 != nil {
		return items, err2
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Access
		err2 := rows.Scan(&i.UUID, &i.Code1, &i.Code2, &i.AccessDate, &i.Zone, &i.Event, &i.SyncDate)
		if err2 != nil {
			return items, err2
		}
		items = append(items, i)
	}

	return items, nil
}

func (r *Repository) SyncAccessUpdateDelete(date time.Time, ok bool) error {
	var sql string
	if ok {
		sql = "DELETE FROM Access WHERE SyncDate = ?"
	} else {
		sql = "UPDATE Access SET SyncDate = NULL WHERE SyncDate = ?"
	}

	stmt, err := r.Db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(date)
	if err != nil {
		return err
	}

	return nil
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
		UUID,
		Code1,
		Code2,
		AccessDate,
		Zone,
		Event
	) values(?, ?, ?, ?, ?, ?)
	`
	stmt, err := r.Db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.UUID, item.Code1, item.Code2, item.AccessDate, item.Zone, item.Event)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) InfoCard(code string) (*models.Card, error) {
	sql := `
	SELECT *
	FROM Cards
	WHERE Code = ?
	`
	stmt, err := r.Db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := r.Db.Query(sql, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var card models.Card
	if rows.Next() {
		err = rows.Scan(&card.Code, &card.DateFrom, &card.DateTo, &card.Enabled, &card.Photo, &card.Deleted)
		if err != nil {
			return nil, err
		}
	}
	return &card, nil
}

/*
func TestCards() []models.Card {
	return []models.Card{
		{Code: "001-42070AED-6E95-4586-91D0-F90BB10D1B7F", DateFrom: time.Now(), DateTo: time.Now().Add(30 * time.Hour).UTC(), Enabled: true, Photo: "", Deleted: false},
		{Code: "001-C343B902-AF36-4A5B-A3A8-C35A11C138E1", DateFrom: time.Now(), DateTo: time.Now().Add(60 * time.Hour).UTC(), Enabled: false, Photo: "", Deleted: false},
		{Code: "001-3F73A575-387C-45F5-9F48-396C5456C268", DateFrom: time.Now(), DateTo: time.Now().Add(60 * time.Hour).UTC(), Enabled: true, Photo: "", Deleted: false},
		{Code: "002-94E5801E-69BE-4268-9B18-F5C4CB7C5187", DateFrom: time.Now(), DateTo: time.Now().Add(60 * time.Hour).UTC(), Enabled: true, Photo: "", Deleted: false},
		{Code: "002-94E5801E-69BE-4268-9B18-F5C4CB7C5181", DateFrom: time.Now(), DateTo: time.Now().Add(60 * time.Hour).UTC(), Enabled: false, Photo: "", Deleted: false},
	}
}
*/
