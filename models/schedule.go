package models

import "time"

type Event struct {
	ID         int
	Title      string
	Speaker    string
	Auditorium string
	Date       time.Time
}
