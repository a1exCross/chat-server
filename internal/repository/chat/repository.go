package chatRepository

import (
	"context"

	"github.com/a1exCross/chat-server/internal/client/db"
	"github.com/a1exCross/chat-server/internal/repository"
)

const tableChats = "chats"

const (
	id        = "id"
	usernames = "usernames"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context) (int64, error) {
	return 0, nil
}
