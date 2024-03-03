package messageRepository

import (
	"context"

	"github.com/a1exCross/chat-server/internal/client/db"
	"github.com/a1exCross/chat-server/internal/repository"
)

const tableMessages = "messages"

const (
	id        = "id"
	author    = "author"
	content   = "content"
	createdAt = "created_at"
	chatID    = "chat_id"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.MessagesRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, message string) error {
	return nil
}
