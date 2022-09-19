package models

import "time"

type Card struct {
	Code    string
	Since   time.Time
	Until   time.Time
	Enabled bool
}

type Access struct {
	Code         string
	AccessDate   time.Time
	Zone         string
	Event        string
	Synchronized bool
}
