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
	UUID       string    `db:"UUID" json:"uuid"`
	Zone       string    `db:"Zone" json:"zone"`
	Event      string    `db:"Event" json:"event"`
	AccessDate time.Time `db:"AccessDate" json:"accessDate"`
	Code1      string    `db:"Code1" json:"code1"`
	Code2      string    `db:"Code2" json:"code2"`
	SyncDate   time.Time `db:"SyncDate" json:"syncDate"`
}

type AccessBulk struct {
	SyncDate time.Time `db:"SyncDate" json:"syncDate"`
	ClientID string    `db:"ClientID" json:"clientId"`
	Access   []Access  `json:"access"`
}
