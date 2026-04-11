package app

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

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
	errCh := make(chan error, 2)

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.diContainer.OrderConsumer(runCtx).RunOrderConsumer(runCtx); err != nil {
			errCh <- fmt.Errorf("order consumer crashed: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.diContainer.AssembledConsumer(runCtx).RunAssembledConsumer(runCtx); err != nil {
			errCh <- fmt.Errorf("assembled consumer crashed: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.runTelegramBot(runCtx)
	}()

	var runErr error

	select {
	case <-runCtx.Done():
		logger.Info(runCtx, "shutdown signal received")
	case err := <-errCh:
		runErr = err
		logger.Error(runCtx, "component crashed, shutting down", zap.Error(err))
	}

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(runCtx, 10*time.Second)
	defer shutdownCancel()

	if closeErr := closer.CloseAll(shutdownCtx); closeErr != nil && !errors.Is(closeErr, context.Canceled) {
		logger.Error(runCtx, "failed to close app resources", zap.Error(closeErr))
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-shutdownCtx.Done():
		logger.Error(runCtx, "shutdown timeout exceeded while waiting for components to stop", zap.Error(shutdownCtx.Err()))
	}

	return runErr
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
