package messageservice

import (
	"context"

	"github.com/a1exCross/chat-server/internal/model"
)

func (s *serv) SendMessage(ctx context.Context, params model.MessageDTO) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		err := s.messageRepo.Create(ctx, params)
		if err != nil {
			return err
		}

		_, err = s.logsRepo.Create(ctx, model.Log{
			Action:  "message created",
			Content: params.Author,
		})

		return err
	})
}
