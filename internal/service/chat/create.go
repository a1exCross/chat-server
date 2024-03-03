package chatservice

import (
	"context"
	"strconv"

	"github.com/a1exCross/chat-server/internal/model"
)

func (s *serv) Create(ctx context.Context, params model.ChatDTO) (int64, error) {
	var id int64

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error

		id, err = s.chatRepo.Create(ctx, params)
		if err != nil {
			return err
		}

		_, err = s.logsRepo.Create(ctx, model.Log{
			Action:  "chat created",
			Content: strconv.FormatInt(id, 10),
		})

		return err
	})

	return id, err
}
