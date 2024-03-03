package chatapi

import (
	"context"
	"fmt"

	"github.com/a1exCross/chat-server/internal/model"
	chatPb "github.com/a1exCross/chat-server/pkg/chat_v1"
)

// Create принимает и обрабатывает запрос на создание чата
func (i *Implementation) Create(ctx context.Context, req *chatPb.CreateRequest) (*chatPb.CreateResponse, error) {
	res, err := i.chatServ.Create(ctx, model.ChatDTO{
		Usernames: req.Usernames,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %v", err)
	}

	return &chatPb.CreateResponse{
		Id: res,
	}, nil
}
