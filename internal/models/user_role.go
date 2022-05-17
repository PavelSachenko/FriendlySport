package models

var UserRoleTable = "users_roles"

type UserRole struct {
	ID          uint64 `json:"id" sql:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
}