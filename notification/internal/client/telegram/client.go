package telegram

import (
	"context"

	"github.com/go-telegram/bot"
)

type client struct {
	bot *bot.Bot
}

func NewTelegramClient(bot *bot.Bot) *client {
	return &client{
		bot: bot,
	}
}

func (c *client) SendMessage(ctx context.Context, chatId int64, message string) error {
	_, err := c.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   message,
	})
	if err != nil {
		return err
	}

	return nil
}
