package tests

import (
	"context"
	"errors"
	"testing"

	chatapi "github.com/a1exCross/chat-server/internal/api/chat"
	"github.com/a1exCross/chat-server/internal/service"
	"github.com/a1exCross/chat-server/internal/service/mocks"
	"github.com/a1exCross/chat-server/pkg/chat_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDelete(t *testing.T) {
	ctx := context.Background()
	mc := minimock.NewController(t)

	type mockChat func(mc *minimock.Controller) service.ChatServive
	type mockMessage func(mc *minimock.Controller) service.MessageService

	id := int64(1)

	deleteReq := &chat_v1.DeleteRequest{
		Id: id,
	}

	tests := []struct {
		name     string
		ctx      context.Context
		err      error
		req      *chat_v1.DeleteRequest
		expected *empty.Empty
		mockMessage
		mockChat
	}{
		{
			name:     "sucessfull test",
			ctx:      ctx,
			err:      nil,
			req:      deleteReq,
			expected: &emptypb.Empty{},
			mockChat: func(mc *minimock.Controller) service.ChatServive {
				mock := mocks.NewChatServiveMock(t)
				mock.DeleteMock.Expect(ctx, id).Return(nil)

				return mock
			},
			mockMessage: func(mc *minimock.Controller) service.MessageService {
				mock := mocks.NewMessageServiceMock(t)

				return mock
			},
		},
		{
			name:     "failed test",
			ctx:      ctx,
			err:      errors.New("failed to delete chat: error"),
			req:      deleteReq,
			expected: nil,
			mockChat: func(mc *minimock.Controller) service.ChatServive {
				mock := mocks.NewChatServiveMock(t)
				mock.DeleteMock.Expect(ctx, id).Return(errors.New("error"))

				return mock
			},
			mockMessage: func(mc *minimock.Controller) service.MessageService {
				mock := mocks.NewMessageServiceMock(t)

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

			res, err := serv.Delete(ctx, test.req)

			require.Equal(t, test.expected, res)
			require.Equal(t, test.err, err)
		})
	}
}
