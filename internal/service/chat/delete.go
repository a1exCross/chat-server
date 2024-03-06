package chatservice

import (
	"context"
	"strconv"

	"github.com/a1exCross/chat-server/internal/model"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		err := s.chatRepo.Delete(ctx, id)
		if err != nil {
			return err
		}

		_, err = s.logsRepo.Create(ctx, model.Log{
			Action:  "chat deleted",
			Content: strconv.FormatInt(id, 10),
		})

		return err
	})
}
