package models

import "fmt"

type SchedulePagination struct {
	isPrev  bool
	Page    int
	MaxPage int
}

func (s SchedulePagination) String() string {
	if s.isPrev {
		return fmt.Sprintf("pager:schedule:prev:%d:%d", s.Page, s.MaxPage)
	}
	return fmt.Sprintf("pager:schedule:next:%d:%d", s.Page, s.MaxPage)
}
