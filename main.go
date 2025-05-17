package main

import (
	"context"
	"github.com/qeery8/api"
	"github.com/qeery8/events/telegram"
	"log"
	"time"
)

func main() {

	ctx := context.Background()

	host := "api.telegram.org"
	token := "TOKEN_BOT"
	client := api.New(host, token)

	processor := telegram.New(client)

	log.Println("Bot started...")

	for {
		events, err := processor.Fetch(ctx, 10)
		if err != nil {
			log.Printf("Fetch error: %v", err)
			continue
		}

		for _, event := range events {
			if err := processor.Process(ctx, event); err != nil {
				log.Printf("Process error: %v", err)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
