package telegram

import (
	"context"
	"errors"
	"github.com/qeery8/api"
	"github.com/qeery8/events"
	e "github.com/qeery8/lib"
	"github.com/qeery8/parsing"
	"log"
)

type Processor struct {
	tg     *api.Client
	offset int
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *api.Client) *Processor {
	return &Processor{
		tg: client,
	}
}

func (p *Processor) Fetch(ctx context.Context, limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(ctx, p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(ctx context.Context, event events.Event) error {
	log.Printf("Processing event of type: %s", event.Type)

	switch event.Type {
	case events.Message:
		return p.processMessage(ctx, event)
	case events.CallbackQuery:
		return p.processCallback(ctx, event)
	default:
		log.Printf("Unknown event type: %v", event.Type)
		return nil
	}
}

func (p *Processor) processCallback(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process callback", err)
	}

	log.Printf("Processing callback with text: %s", event.Text)

	switch event.Text {
	case "site_wallapop":
		return p.handleWallapop(ctx, meta.ChatID)
	//case "site_subito":
	//return p.handleSubito(ctx, meta.ChatID)
	//case "site_fiverr":
	//return p.HandleFiverr(ctx, meta.ChatID)
	default:
		return p.tg.SendMessage(ctx, meta.ChatID, "Unknown command")
	}
}

func (p *Processor) handleWallapop(ctx context.Context, chatID int) error {
	data, err := parsing.ParseWallapop()
	if err != nil {
		return p.tg.SendMessage(ctx, chatID, "Failed to fetch Wallapop ads.")
	}
	return p.tg.SendMessage(ctx, chatID, data)
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(ctx, event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(upd api.Update) events.Event {
	updType := fetchType(upd)

	log.Printf("Fetched event type: %v", updType)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	if updType == events.CallbackQuery {
		if upd.CallbackQuery != nil {
			res.Meta = Meta{
				ChatID:   upd.CallbackQuery.Message.Chat.ID,
				Username: upd.CallbackQuery.From.Username,
			}
			res.Text = upd.CallbackQuery.Data
		}
	}

	return res
}

func fetchText(upd api.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd api.Update) events.Type {
	if upd.Message != nil && upd.Message.Text != "" {
		log.Printf("Message received: %v", upd.Message)
		return events.Message
	}

	if upd.CallbackQuery != nil {
		log.Printf("CallbackQuery received: %v", upd.CallbackQuery)
		return events.CallbackQuery
	}

	log.Printf("Unknown event: %v", upd)
	return events.Unknown
}
