package pgxtx

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/ChopX4/raketka/platform/pkg/logger"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type txManager struct {
	db *pgxpool.Pool
}

type txCtxKey struct{}

func NewTxManager(db *pgxpool.Pool) TxManager {
	return &txManager{
		db: db,
	}
}

func (m *txManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	if _, ok := TxFromContext(ctx); ok {
		return fn(ctx)
	}

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
				logger.Error(ctx, "failed to rollback transaction", zap.Error(rollbackErr))
			}
			return
		}

		err = tx.Commit(ctx)
	}()

	return fn(ContextWithTx(ctx, tx))
}

func ContextWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txCtxKey{}, tx)
}

func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txCtxKey{}).(pgx.Tx)
	return tx, ok
}
