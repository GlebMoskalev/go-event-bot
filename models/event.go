package models

import "time"

type Status string

const (
	StatusPlanned   Status = "planned"
	StatusOngoing   Status = "ongoing"
	StatusCompleted Status = "completed"
)

type Event struct {
	ID         int
	Title      string
	Speaker    string
	Auditorium string
	Status     Status
	Date       time.Time
}
