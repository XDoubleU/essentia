package database

import (
	"context"
	"sync"
)

// MinimalDBTx is the minimal interface a DB transaction should implement.
type MinimalDBTx interface {
	Rollback(ctx context.Context) error
}

// SyncTx wraps a database transaction to make sure it can be used concurrently.
// This is achieved by locking the transaction when in use.
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

// CreateSyncTx creates a [SyncTx] which
// makes sure a database transaction can be used concurrently.
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

// WrapInSyncTx is used to make sure a
// transactional database action can run concurrently.
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

// WrapInSyncTxNoError is used to make sure a
// transactional database action can run concurrently.
// The executed database action shouldn't return an error.
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

// Rollback is used to rollback the wrapped transaction.
func (tx SyncTx[TTx]) Rollback(ctx context.Context) error {
	waitOnLock(tx.mu)
	defer tx.mu.Unlock()

	return tx.Tx.Rollback(ctx)
}
