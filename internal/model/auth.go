package model

// AccessCheck - хранит роли пользователей
type AccessCheck struct {
	Roles []UserRole
}

// UserRole - Описывает роль пользователя от 0 до 255
type UserRole int8

// Типы ролей пользователя
const (
	UNKNOWN UserRole = iota
	USER
	ADMIN
)
