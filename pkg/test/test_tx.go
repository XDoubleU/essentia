package test

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Tx struct {
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

func (tx Tx) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (pgconn.CommandTag, error) {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	return tx.tx.Exec(ctx, sql, arguments...)
}

func (tx Tx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	//nolint:sqlclosecheck // user is supposed to close query
	return tx.tx.Query(ctx, sql, args...)
}

func (tx Tx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	return tx.tx.QueryRow(ctx, sql, args...)
}

func (tx Tx) Begin(ctx context.Context) (pgx.Tx, error) {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	return tx.tx.Begin(ctx)
}
