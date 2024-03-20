package chatapi

import (
	"context"
	"fmt"

	"github.com/a1exCross/chat-server/internal/api/chat/converter"
	chatPb "github.com/a1exCross/chat-server/pkg/chat_v1"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SendMessage принимает и обрабатывает запрос на отправку сообщения в чате
func (i *Implementation) SendMessage(ctx context.Context, req *chatPb.SendMessageRequest) (*empty.Empty, error) {
	err := i.messageServ.SendMessage(ctx, converter.ProtoToMessage(req.Message))
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %v", err)
	}

	return &emptypb.Empty{}, nil
}
