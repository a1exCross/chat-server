package chatAPI

import (
	"context"

	chatPb "github.com/a1exCross/chat-server/pkg/chat_v1"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Create(ctx context.Context, req *chatPb.CreateRequest) (*chatPb.CreateResponse, error) {
	/* insertBuilder := sq.Insert(tableChats).
		PlaceholderFormat(sq.Dollar).
		Columns(usernames).
		Values(req.Usernames).
		Suffix(fmt.Sprintf("RETURNING %s", id))

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error at parse sql builder: %v", err)
	}

	var id int64

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error at query to database: %v", err)
	}
	*/
	return &chatPb.CreateResponse{
		Id: 0,
	}, nil
}

func (i *Implementation) Delete(ctx context.Context, req *chatPb.DeleteRequest) (*empty.Empty, error) {
	/* deleteBuilder := sq.Delete(tableChats).
		Where(fmt.Sprintf("%s = ?", id), req.GetId()).
		PlaceholderFormat(sq.Dollar)

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error at parse sql builder: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error at query to database: %v", err)
	}

	deleteBuilder = sq.Delete(tableMessages).
		Where(fmt.Sprintf("%s = ?", chatID), req.GetId()).
		PlaceholderFormat(sq.Dollar)

	query, args, err = deleteBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error at parse sql builder: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error at query to database: %v", err)
	} */

	return &empty.Empty{}, nil
}

func (i *Implementation) SendMessage(ctx context.Context, req *chatPb.SendMessageRequest) (*empty.Empty, error) {
	/* query, args, err := sq.Select(id).
		From(tableChats).
		OrderBy("id DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error at parse sql builder: %v", err)
	}

	var lastChatID int64

	err = s.pool.QueryRow(ctx, query, args...).Scan(&lastChatID)
	if err != nil {
		return nil, fmt.Errorf("error at query to database: %v", err)
	}

	insertBuilder := sq.Insert(tableMessages).
		PlaceholderFormat(sq.Dollar).
		Columns(chatID, author, content, createdAt).
		Values(lastChatID, req.Message.From, req.Message.Text, req.Message.Timestamp.AsTime()).
		Suffix(fmt.Sprintf("RETURNING %s", id))

	query, args, err = insertBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error at parse sql builder: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error at query to database: %v", err)
	}
	*/
	return &emptypb.Empty{}, nil
}
