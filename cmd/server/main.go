package main

import (
	"context"
	"log"

	"github.com/a1exCross/chat-server/internal/app"
)

func main() {

	ctx := context.Background()

	app, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
