package models

import (
	"fmt"
	"strings"
)

type Direction string

const (
	Next Direction = "Next"
	Prev Direction = "Prev"

	PaginationPrefix = "pagination"
	EventContext     = "event"
	AllContext       = "all"
	DaysContext      = "dayas"

	ItemsPerPage = 5
)

type CallbackButton struct {
	Text string
	Data string
}

func PaginationEvent(currentPage, maxPage int, direction Direction) CallbackButton {
	btn := CallbackButton{}
	btn.Text = string(direction)
	btn.Data = fmt.Sprintf(
		"%s:%s:%s:%d:%d",
		PaginationPrefix,
		EventContext,
		strings.ToLower(string(direction)), currentPage, maxPage,
	)
	return btn
}

func PageNumber(currentPage, maxPage int) CallbackButton {
	btn := CallbackButton{}
	btn.Text = fmt.Sprintf("%d / % d", currentPage, maxPage)
	btn.Data = "null"
	return btn
}
