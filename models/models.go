package models

import "time"

type Card struct {
	Code    string
	Since   time.Time
	Until   time.Time
	Enabled bool
}

type Access struct {
	Code1      string    `db:"Code1" json:"code1"`
	Code2      string    `db:"Code2" json:"code2"`
	AccessDate time.Time `db:"AccessDate" json:"accessDate"`
	Zone       string    `db:"Zone" json:"zone"`
	Event      string    `db:"Event" json:"event"`
}

type AccessZone struct {
	ClientID string   `db:"ClientID" json:"clientId"`
	Access   []Access `json:"access"`
}
