package models

type Role string

const (
	RoleGuest = "guest"
	RoleStaff = "staff"
	RoleAdmin = "admin"
)

type User struct {
	TelegramID int64
	FirstName  string
	LastName   string
	Patronymic string
	Role       Role
}

func (u *User) HasRole(role Role) bool {
	switch role {
	case RoleGuest:
		return true
	case RoleStaff:
		return u.Role == RoleStaff || u.Role == RoleAdmin
	case RoleAdmin:
		return u.Role == RoleAdmin
	default:
		return false
	}
}
