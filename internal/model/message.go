package model

import "time"

type Message struct {
	ID        int64     `db:"id"`
	ChatID    int64     `db:"chat_id"`
	Author    string    `db:"author"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

type MessageDTO struct {
	Author  string `db:"author"`
	Content string `db:"content"`
}
