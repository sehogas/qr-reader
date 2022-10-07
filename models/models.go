package models

import "time"

type Card struct {
	Code     string    `db:"Code" json:"code"`
	DateFrom time.Time `db:"DateFrom" json:"from"`
	DateTo   time.Time `db:"DateTo" json:"to"`
	Enabled  bool      `db:"Enabled" json:"enabled"`
	Photo    string    `db:"Photo" json:"photo"`
	Deleted  bool      `db:"Deleted" json:"deleted"`
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
