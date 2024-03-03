package model

import "time"

// Message - структура, оисывающая сообщение в БД
type Message struct {
	ID        int64     `db:"id"`
	ChatID    int64     `db:"chat_id"`
	Author    string    `db:"author"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

// MessageDTO - DTO для создания нового сообщения
type MessageDTO struct {
	Author    string
	Content   string
	CreatedAt time.Time
}
