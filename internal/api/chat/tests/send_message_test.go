package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	chatapi "github.com/a1exCross/chat-server/internal/api/chat"
	"github.com/a1exCross/chat-server/internal/model"
	"github.com/a1exCross/chat-server/internal/service"
	"github.com/a1exCross/chat-server/internal/service/mocks"
	"github.com/a1exCross/chat-server/pkg/chat_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSend(t *testing.T) {
	ctx := context.Background()
	mc := minimock.NewController(t)

	type mockChat func(mc *minimock.Controller) service.ChatServive
	type mockMessage func(mc *minimock.Controller) service.MessageService

	timeNow := time.Now().UTC()

	messageDTO := model.MessageDTO{
		Author:    "biba",
		Content:   "hello",
		CreatedAt: timeNow,
	}

	createReq := &chat_v1.SendMessageRequest{
		Message: &chat_v1.Message{
			From:      "biba",
			Text:      "hello",
			Timestamp: timestamppb.New(timeNow),
		},
	}

	tests := []struct {
		name     string
		ctx      context.Context
		err      error
		req      *chat_v1.SendMessageRequest
		expected *empty.Empty
		mockMessage
		mockChat
	}{
		{
			name:     "sucessfull test",
			ctx:      ctx,
			err:      nil,
			req:      createReq,
			expected: &emptypb.Empty{},
			mockChat: func(mc *minimock.Controller) service.ChatServive {
				mock := mocks.NewChatServiveMock(t)

				return mock
			},
			mockMessage: func(mc *minimock.Controller) service.MessageService {
				mock := mocks.NewMessageServiceMock(t)
				mock.SendMessageMock.Expect(ctx, messageDTO).Return(nil)

				return mock
			},
		},
		{
			name:     "failed test",
			ctx:      ctx,
			err:      errors.New("failed to send message: error"),
			req:      createReq,
			expected: nil,
			mockChat: func(mc *minimock.Controller) service.ChatServive {
				mock := mocks.NewChatServiveMock(t)

				return mock
			},
			mockMessage: func(mc *minimock.Controller) service.MessageService {
				mock := mocks.NewMessageServiceMock(t)
				mock.SendMessageMock.Expect(ctx, messageDTO).Return(errors.New("error"))

				return mock
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			chatService := test.mockChat(mc)
			messageService := test.mockMessage(mc)

			serv := chatapi.NewImplementation(chatService, messageService)

			res, err := serv.SendMessage(ctx, test.req)

			require.Equal(t, test.expected, res)
			require.Equal(t, test.err, err)
		})
	}
}
