package postgres

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SyncTx struct {
	tx pgx.Tx
	mu *sync.Mutex
}

func waitOnLock(lock *sync.Mutex) {
	for {
		if lock.TryLock() {
			break
		}
	}
}

func CreateSyncTx(ctx context.Context, pool *pgxpool.Pool) SyncTx {
	var mu sync.Mutex
	for {
		tx, err := pool.Begin(ctx)
		if err == nil {
			return SyncTx{
				tx: tx,
				mu: &mu,
			}
		}
	}
}

func (tx SyncTx) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (pgconn.CommandTag, error) {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	return tx.tx.Exec(ctx, sql, arguments...)
}

func (tx SyncTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	//nolint:sqlclosecheck // user is supposed to close query
	return tx.tx.Query(ctx, sql, args...)
}

func (tx SyncTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	return tx.tx.QueryRow(ctx, sql, args...)
}

func (tx SyncTx) Begin(ctx context.Context) (pgx.Tx, error) {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	return tx.tx.Begin(ctx)
}

func (tx SyncTx) Rollback(ctx context.Context) error {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	return tx.tx.Rollback(ctx)
}
