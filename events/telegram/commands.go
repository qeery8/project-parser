package telegram

import (
	"context"
	"log"
	"strings"
)

const (
	StartCmd = "/start"
	HelpCmd  = "/help"
	ParsCmd  = "/pars"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	switch text {
	case StartCmd:
		return p.sendStart(ctx, chatID)
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case ParsCmd:
		return p.sendSitChoice(ctx, chatID)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) sendSitChoice(ctx context.Context, chatID int) error {
	msg := "the site from which you want to receive the ad"
	keyboard := `{
		"inline_keyboard": [
			[
				{"text": "Wallapop", "callback_data": "site_wallapop"},
				{"text": "Subito", "callback_data": "site_subito"},
				{"text": "Fiverr", "callback_data": "site_fiverr"}
			]
		]
	}`

	return p.tg.SendMessage(ctx, chatID)
}

func (p *Processor) sendStart(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgStart)
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}
