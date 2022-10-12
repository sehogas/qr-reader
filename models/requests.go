package models

import "time"

type QRUpdatesResponse struct {
	ServerTime time.Time `db:"ServerTime" json:"serverTime"`
	Cards      []Card    `db:"Cards" json:"cards"`
}

type AccessDataResponse struct {
	DocumentType   string    `json:"docType"`
	DocumentNumber string    `json:"docNum"`
	DateFrom       time.Time `json:"dateFrom"`
	DateTo         time.Time `json:"dateTo"`
	Color          string    `json:"color"`
	PersonName     string    `json:"personName"`
	PNA            string    `json:"pna"`
	Photo          string    `json:"photo"`
	Eventual       bool      `json:"eventual"`
	LicensePlate   string    `json:"licensePlate"`
}
