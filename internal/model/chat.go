package model

type Chat struct {
	ID        int64    `db:"id"`
	Usernames []string `db:"usernames"`
}

type ChatDTO struct {
	Usernames []string `db:"usernames"`
}
