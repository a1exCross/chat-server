package model

// Chat - структура, описывающая чат в БД
type Chat struct {
	ID        int64    `db:"id"`
	Usernames []string `db:"usernames"`
}

// ChatDTO - DTO для создания нового чата
type ChatDTO struct {
	Usernames []string
}
