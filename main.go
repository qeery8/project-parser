package main

import (
	"context"
	"github.com/qeery8/api"
)

func main() {
	ctx := context.Background()

	host := "api_telegram"
	token := "TOKEN_BOT"
	client := api.New(host, token)
}
