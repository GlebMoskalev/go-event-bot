package apperrors

import "errors"

var (
	ErrNotFoundStaff      = errors.New("not found staff")
	ErrInvalidPhoneNumber = errors.New("invalid phone number format")
	ErrNotFoundUser       = errors.New("user not found")
	ErrNotFoundSchedule   = errors.New("user not found")
)
