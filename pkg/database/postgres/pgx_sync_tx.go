package postgres

import (
	"context"

	"github.com/XDoubleU/essentia/pkg/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgxSyncTx uses [database.SyncTx] to make sure
// [pgx.Tx] can be used concurrently.
type PgxSyncTx struct {
	syncTx database.SyncTx[pgx.Tx]
}

// CreatePgxSyncTx returns a [pgx.Tx] which works concurrently.
func CreatePgxSyncTx(ctx context.Context, db *pgxpool.Pool) PgxSyncTx {
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
	return database.WrapInSyncTx(ctx, tx.syncTx, tx.syncTx.Tx.Exec, sql, args...)
}

// Query is used to wrap [pgx.Tx.Query] in a [database.SyncTx].
func (tx PgxSyncTx) Query(
	ctx context.Context,
	sql string,
	args ...any,
) (pgx.Rows, error) {
	return database.WrapInSyncTx(ctx, tx.syncTx, tx.syncTx.Tx.Query, sql, args...)
}

// QueryRow is used to wrap [pgx.Tx.QueryRow] in a [database.SyncTx].
func (tx PgxSyncTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return database.WrapInSyncTxNoError(
		ctx,
		tx.syncTx,
		tx.syncTx.Tx.QueryRow,
		sql,
		args...)
}

// Rollback is used to rollback a [PgxSyncTx].
func (tx PgxSyncTx) Rollback(ctx context.Context) error {
	return tx.syncTx.Rollback(ctx)
}
