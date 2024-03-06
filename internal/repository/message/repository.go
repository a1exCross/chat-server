package messagerepository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/a1exCross/chat-server/internal/client/db"
	"github.com/a1exCross/chat-server/internal/model"
	"github.com/a1exCross/chat-server/internal/repository"
)

const (
	chatsTable    = "chats"
	messagesTable = "messages"
	logsTable     = "logs"

	idColumn        = "id"
	contentColumn   = "content"
	authorColumn    = "author"
	createdAtColumn = "created_at"
	chatIDColumn    = "chat_id"
)

type repo struct {
	db db.Client
}

// NewRepository - возвращает методы для работы с репозиторием сообщений
func NewRepository(db db.Client) repository.MessagesRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, params model.MessageDTO) error {
	query, args, err := sq.Select(idColumn).
		From(chatsTable).
		OrderBy(fmt.Sprintf("%s DESC", idColumn)).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("error at parse sql builder: %v", err)
	}

	var lastChatID int64

	q := db.Query{
		Name:     "messages_repository.Create.GetLastChat",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&lastChatID)
	if err != nil {
		return fmt.Errorf("error at query to database: %v", err)
	}

	insertBuilder := sq.Insert(messagesTable).
		PlaceholderFormat(sq.Dollar).
		Columns(chatIDColumn, authorColumn, contentColumn, createdAtColumn).
		Values(lastChatID, params.Author, params.Content, params.CreatedAt).
		Suffix(fmt.Sprintf("RETURNING %s", idColumn))

	query, args, err = insertBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error at parse sql builder: %v", err)
	}

	q = db.Query{
		Name:     "messages_repository.Create.CreateMessage",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("error at query to database: %v", err)
	}

	return nil
}
