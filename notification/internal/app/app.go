package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/ChopX4/raketka/platform/pkg/closer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 3)

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.diContainer.OrderConsumer(runCtx).RunOrderConsumer(runCtx); err != nil {
			errCh <- fmt.Errorf("order consumer crashed: %w", err)
		}
	}()

	go func() {
		if err := a.diContainer.AssembledConsumer(runCtx).RunAssembledConsumer(runCtx); err != nil {
			errCh <- fmt.Errorf("assembled consumer crashed: %w", err)
		}
	}()

	go func() {
		a.runTelegramBot(runCtx)
	}()

	select {
	case <-runCtx.Done():
		return runCtx.Err()
	case err := <-errCh:
		logger.Error(runCtx, "component crashed, shutting down", zap.Error(err))
		cancel()

		shutdownCtx, shutdownCancel := context.WithCancel(runCtx)
		defer shutdownCancel()

		if closeErr := closer.CloseAll(shutdownCtx); closeErr != nil && !errors.Is(closeErr, context.Canceled) {
			logger.Error(runCtx, "failed to close app resources", zap.Error(closeErr))
		}

		return err
	}
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initTelegramBot,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) initTelegramBot(ctx context.Context) error {
	telegramBot := a.diContainer.TelegramBot(ctx)

	telegramBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		logger.Info(ctx, "telegram /start received", zap.Int64("chat_id", update.Message.Chat.ID))

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Notification bot activated. You will receive order updates here.",
		})
		if err != nil {
			logger.Error(ctx, "failed to send activation message", zap.Error(err))
		}
	})

	return nil
}

func (a *App) runTelegramBot(ctx context.Context) {
	logger.Info(ctx, "telegram bot started")
	a.diContainer.TelegramBot(ctx).Start(ctx)
}
