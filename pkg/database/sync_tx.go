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
	sync.Mutex
	Tx TTx
}

// CreateSyncTx creates a [SyncTx] which
// makes sure a database transaction can be used concurrently.
func CreateSyncTx[TTx MinimalDBTx](
	ctx context.Context,
	beginTxFunc func(ctx context.Context) (TTx, error),
) *SyncTx[TTx] {
	for {
		tx, err := beginTxFunc(ctx)
		if err == nil {
			return &SyncTx[TTx]{
				Tx: tx,
			}
		}
	}
}

// WrapInSyncTx is used to make sure a
// transactional database action can run concurrently.
func WrapInSyncTx[TTx MinimalDBTx, TResult any](
	ctx context.Context,
	tx *SyncTx[TTx],
	queryFunc func(ctx context.Context) (TResult, error),
) (TResult, error) {
	tx.Lock()
	defer tx.Unlock()

	return queryFunc(ctx)
}

// WrapInSyncTxNoError is used to make sure a
// transactional database action can run concurrently.
// The executed database action shouldn't return an error.
func WrapInSyncTxNoError[TTx MinimalDBTx, TResult any](
	ctx context.Context,
	tx *SyncTx[TTx],
	queryFunc func(ctx context.Context) TResult,
) TResult {
	tx.Lock()
	defer tx.Unlock()

	return queryFunc(ctx)
}

// WrapInSyncTxNoReturn is used to make sure a
// transactional database action can run concurrently.
// The executed database action shouldn't return anything.
func WrapInSyncTxNoReturn[TTx MinimalDBTx](
	ctx context.Context,
	tx *SyncTx[TTx],
	queryFunc func(ctx context.Context),
) {
	tx.Lock()
	defer tx.Unlock()

	queryFunc(ctx)
}

// Rollback is used to rollback the wrapped transaction.
func (tx *SyncTx[TTx]) Rollback(ctx context.Context) error {
	tx.Lock()
	defer tx.Unlock()

	return tx.Tx.Rollback(ctx)
}
