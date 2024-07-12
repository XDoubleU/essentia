package postgres

import (
	"context"

	"github.com/XDoubleU/essentia/pkg/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxSyncTx struct {
	syncTx database.SyncTx[pgx.Tx]
}

func CreatePgxSyncTx(ctx context.Context, db *pgxpool.Pool) PgxSyncTx {
	syncTx := database.CreateSyncTx(ctx, db.Begin)

	return PgxSyncTx{
		syncTx: syncTx,
	}
}

func (tx PgxSyncTx) Exec(
	ctx context.Context,
	sql string,
	args ...any,
) (pgconn.CommandTag, error) {
	return database.WrapInSyncTx(ctx, tx.syncTx, tx.syncTx.Tx.Exec, sql, args...)
}

func (tx PgxSyncTx) Query(
	ctx context.Context,
	sql string,
	args ...any,
) (pgx.Rows, error) {
	return database.WrapInSyncTx(ctx, tx.syncTx, tx.syncTx.Tx.Query, sql, args...)
}

func (tx PgxSyncTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return database.WrapInSyncTxNoError(
		ctx,
		tx.syncTx,
		tx.syncTx.Tx.QueryRow,
		sql,
		args...)
}

func (tx PgxSyncTx) Rollback(ctx context.Context) error {
	return tx.syncTx.Rollback(ctx)
}
