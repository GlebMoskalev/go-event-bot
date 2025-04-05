package apperrors

import "errors"

var (
	ErrNotFoundStaff      = errors.New("not found staff")
	ErrInvalidPhoneNumber = errors.New("invalid phone number format")
	ErrNotFoundUser       = errors.New("user not found")
	ErrNotFoundEvent      = errors.New("user not found")
	ErrUserExists         = errors.New("user exists")
	ErrNotFoundState      = errors.New("not found state")
	ErrIncompleteFullName = errors.New("name is incomplete")
)
