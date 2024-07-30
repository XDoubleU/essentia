package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xdoubleu/essentia/pkg/database"
)

// PgxSyncTx uses [database.SyncTx] to make sure
// [pgx.Tx] can be used concurrently.
type PgxSyncTx struct {
	syncTx database.SyncTx[pgx.Tx]
}

// CreatePgxSyncTx returns a [pgx.Tx] which works concurrently.
func CreatePgxSyncTx(ctx context.Context, db DB) PgxSyncTx {
	syncTx := database.CreateSyncTx(ctx, db.Begin)

	return PgxSyncTx{
		syncTx: syncTx,
	}
}

// Exec is used to wrap [pgx.Tx.Exec] in a [database.SyncTx].
func (tx PgxSyncTx) Exec(
	ctx context.Context,
	sql string,
	args ...any,
) (pgconn.CommandTag, error) {
	return database.WrapInSyncTx(ctx, tx.syncTx, func(ctx context.Context) (pgconn.CommandTag, error) {
		return tx.syncTx.Tx.Exec(ctx, sql, args...)
	})
}

// Query is used to wrap [pgx.Tx.Query] in a [database.SyncTx].
func (tx PgxSyncTx) Query(
	ctx context.Context,
	sql string,
	args ...any,
) (pgx.Rows, error) {
	return database.WrapInSyncTx(ctx, tx.syncTx, func(ctx context.Context) (pgx.Rows, error) {
		return tx.syncTx.Tx.Query(ctx, sql, args...)
	})
}

// QueryRow is used to wrap [pgx.Tx.QueryRow] in a [database.SyncTx].
func (tx PgxSyncTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return database.WrapInSyncTxNoError(ctx, tx.syncTx, func(ctx context.Context) pgx.Row {
		return tx.syncTx.Tx.QueryRow(ctx, sql, args...)
	})
}

// Begin is used to begin a [pgx.Tx].
func (tx PgxSyncTx) Begin(ctx context.Context) (pgx.Tx, error) {
	return database.WrapInSyncTx(ctx, tx.syncTx, func(ctx context.Context) (pgx.Tx, error) {
		return tx.syncTx.Tx.Begin(ctx)
	})
}

func (tx PgxSyncTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return database.WrapInSyncTxNoError(ctx, tx.syncTx, func(ctx context.Context) pgx.BatchResults {
		return tx.syncTx.Tx.SendBatch(ctx, b)
	})
}

// Rollback is used to rollback a [PgxSyncTx].
func (tx PgxSyncTx) Rollback(ctx context.Context) error {
	return tx.syncTx.Rollback(ctx)
}
