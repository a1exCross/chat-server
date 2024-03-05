package chatrepository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/a1exCross/chat-server/internal/model"
	"github.com/a1exCross/chat-server/internal/repository"
	"github.com/a1exCross/common/pkg/client/db"
)

const (
	logsTable     = "logs"
	chatsTable    = "chats"
	messagesTable = "messages"

	idColumn        = "id"
	contentColumn   = "content"
	createdAtColumn = "created_at"
	actionColumn    = "action"
	usernamesColumn = "usernames"
	authorColumn    = "author"
	chatIDColumn    = "chat_id"
)

type repo struct {
	db db.Client
}

// NewRepository - возвращает методы для работы с репозиторием чатов
func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, params model.ChatDTO) (int64, error) {
	insertBuilder := sq.Insert(chatsTable).
		PlaceholderFormat(sq.Dollar).
		Columns(usernamesColumn).
		Values(params.Usernames).
		Suffix(fmt.Sprintf("RETURNING %s", idColumn))

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("error at parse sql builder: %v", err)
	}

	q := db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: query,
	}

	var id int64

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error at query to database: %v", err)
	}

	return id, nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	deleteBuilder := sq.Delete(chatsTable).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error at parse sql builder: %v", err)
	}

	q := db.Query{
		Name:     "chat_repository.Delete.Chats",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("error at query to database: %v", err)
	}

	deleteBuilder = sq.Delete(messagesTable).
		Where(sq.Eq{chatIDColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err = deleteBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error at parse sql builder: %v", err)
	}

	q = db.Query{
		Name:     "chat_repository.Delete.Messages",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("error at query to database: %v", err)
	}

	return nil
}
