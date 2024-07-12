package database

import (
	"context"
	"sync"
)

type MinimalDBTx interface {
	Rollback(ctx context.Context) error
}

type SyncTx[TTx MinimalDBTx] struct {
	Tx TTx
	mu *sync.Mutex
}

func waitOnLock(lock *sync.Mutex) {
	for {
		if lock.TryLock() {
			break
		}
	}
}

func CreateSyncTx[TTx MinimalDBTx](
	ctx context.Context,
	beginTxFunc func(ctx context.Context) (TTx, error),
) SyncTx[TTx] {
	var mu sync.Mutex
	for {
		tx, err := beginTxFunc(ctx)
		if err == nil {
			return SyncTx[TTx]{
				Tx: tx,
				mu: &mu,
			}
		}
	}
}

func WrapInSyncTx[TTx MinimalDBTx, TResult any](
	ctx context.Context,
	tx SyncTx[TTx],
	queryFunc func(ctx context.Context, sql string, args ...any) (TResult, error),
	sql string,
	args ...any,
) (TResult, error) {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	return queryFunc(ctx, sql, args...)
}

func WrapInSyncTxNoError[TTx MinimalDBTx, TResult any](
	ctx context.Context,
	tx SyncTx[TTx],
	queryFunc func(ctx context.Context, sql string, args ...any) TResult,
	sql string,
	args ...any,
) TResult {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	return queryFunc(ctx, sql, args...)
}

func (tx SyncTx[TTx]) Rollback(ctx context.Context) error {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	return tx.Tx.Rollback(ctx)
}
