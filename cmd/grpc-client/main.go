package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	descChat "github.com/a1exCross/chat-server/pkg/chat_v1"
)

const (
	grpcPort = 50052
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connection grpc server: %v", err)
	}

	cl := descChat.NewChatV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := cl.Create(ctx, &descChat.CreateRequest{
		Usernames: []string{"biba", "boba"},
	})
	if err != nil {
		log.Fatalf("failed to create chat: %v", err)
	}

	_, err = cl.Delete(ctx, &descChat.DeleteRequest{
		Id: res.Id,
	})
	if err != nil {
		log.Fatalf("failed to delete chat: %v", err)
	}

	_, err = cl.SendMessage(ctx, &descChat.SendMessageRequest{
		Message: &descChat.Message{
			From:      "biba",
			Text:      "zdarova",
			Timestamp: timestamppb.Now(),
		},
	})
	if err != nil {
		log.Fatalf("failed to send message: %v", err)
	}

}
