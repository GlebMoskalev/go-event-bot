package models

import "time"

type Schedule struct {
	ID          int
	Title       string
	Description string
	Date        time.Time
}
