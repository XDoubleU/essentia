package postgres

import (
	"context"

	"github.com/XDoubleU/essentia/pkg/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// PgxSyncTx uses [database.SyncTx] to make sure
// [pgx.Tx] can be used concurrently.
type PgxSyncTx struct {
	syncTx *database.SyncTx[pgx.Tx]
}

// PgxSyncRow is a concurrent wrapper for [pgx.Row].
type PgxSyncRow struct {
	syncTx *database.SyncTx[pgx.Tx]
	row    pgx.Row
}

// PgxSyncRows is a concurrent wrapper for [pgx.Rows].
type PgxSyncRows struct {
	syncTx *database.SyncTx[pgx.Tx]
	rows   pgx.Rows
}

// CreatePgxSyncTx returns a [pgx.Tx] which works concurrently.
func CreatePgxSyncTx(ctx context.Context, db DB) *PgxSyncTx {
	syncTx := database.CreateSyncTx(ctx, db.Begin)

	return &PgxSyncTx{
		syncTx: syncTx,
	}
}

// Exec is used to wrap [pgx.Tx.Exec] in a [database.SyncTx].
func (tx *PgxSyncTx) Exec(
	ctx context.Context,
	sql string,
	args ...any,
) (pgconn.CommandTag, error) {
	return database.WrapInSyncTx(
		ctx,
		tx.syncTx,
		func(ctx context.Context) (pgconn.CommandTag, error) {
			return tx.syncTx.Tx.Exec(ctx, sql, args...)
		},
	)
}

// Query is used to wrap [pgx.Tx.Query] in a [database.SyncTx].
func (tx *PgxSyncTx) Query(
	ctx context.Context,
	sql string,
	args ...any,
) (pgx.Rows, error) {
	tx.syncTx.Mutex.Lock()

	rows, err := tx.syncTx.Tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return &PgxSyncRows{
		syncTx: tx.syncTx,
		rows:   rows,
	}, nil
}

// SendBatch is used to wrap [pgx.Tx.QueryRow] in a [database.SyncTx].
func (tx *PgxSyncTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return database.WrapInSyncTxNoError(
		ctx,
		tx.syncTx,
		func(ctx context.Context) pgx.BatchResults {
			return tx.syncTx.Tx.SendBatch(ctx, b)
		},
	)
}

// Close closes the opened [pgx.Rows].
func (rows *PgxSyncRows) Close() {
	rows.syncTx.Unlock()
	rows.rows.Close()
}

// CommandTag fetches the [pgconn.CommandTag].
func (rows *PgxSyncRows) CommandTag() pgconn.CommandTag {
	return rows.rows.CommandTag()
}

// Conn fetches the [pgx.Conn].
func (rows *PgxSyncRows) Conn() *pgx.Conn {
	return rows.rows.Conn()
}

// Err fetches any errors.
func (rows *PgxSyncRows) Err() error {
	return rows.rows.Err()
}

// FieldDescriptions fetches [pgconn.FieldDescription]s.
func (rows *PgxSyncRows) FieldDescriptions() []pgconn.FieldDescription {
	return rows.rows.FieldDescriptions()
}

// Next continues to the next row of [PgxSyncRows] if there is one.
func (rows *PgxSyncRows) Next() bool {
	return rows.rows.Next()
}

// RawValues fetches the raw values of the current row.
func (rows *PgxSyncRows) RawValues() [][]byte {
	return rows.rows.RawValues()
}

// Scan scans the data of the current row into dest.
func (rows *PgxSyncRows) Scan(dest ...any) error {
	return rows.rows.Scan(dest...)
}

// Values fetches the values of the current row.
func (rows *PgxSyncRows) Values() ([]any, error) {
	return rows.rows.Values()
}

// QueryRow is used to wrap [pgx.Tx.QueryRow] in a [database.SyncTx].
func (tx *PgxSyncTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	tx.syncTx.Mutex.Lock()

	row := tx.syncTx.Tx.QueryRow(ctx, sql, args...)

	return &PgxSyncRow{
		syncTx: tx.syncTx,
		row:    row,
	}
}

// Scan scans the data of [PgxSyncRow] into dest.
func (row *PgxSyncRow) Scan(dest ...any) error {
	defer row.syncTx.Unlock()
	return row.row.Scan(dest...)
}

// Ping is used to wrap [pgx.Tx.Conn.Ping] in a [database.SyncTx].
func (tx *PgxSyncTx) Ping(ctx context.Context) error {
	return database.WrapInSyncTxNoError(
		ctx,
		tx.syncTx,
		func(ctx context.Context) error {
			return tx.syncTx.Tx.Conn().Ping(ctx)
		},
	)
}

// Begin is used to wrap [pgx.Tx.Begin] in a [database.SyncTx].
func (tx *PgxSyncTx) Begin(ctx context.Context) (pgx.Tx, error) {
	return database.WrapInSyncTx(
		ctx,
		tx.syncTx,
		func(ctx context.Context) (pgx.Tx, error) {
			return tx.syncTx.Tx.Begin(ctx)
		},
	)
}

// BeginTx is used to wrap [pgx.Tx.BeginTx] in a [database.SyncTx].
func (tx *PgxSyncTx) BeginTx(
	ctx context.Context,
	txOptions pgx.TxOptions,
) (pgx.Tx, error) {
	return database.WrapInSyncTx(
		ctx,
		tx.syncTx,
		func(ctx context.Context) (pgx.Tx, error) {
			return tx.syncTx.Tx.Conn().BeginTx(ctx, txOptions)
		},
	)
}

// Commit is used to wrap [pgx.Tx.Commit] in a [database.SyncTx].
func (tx *PgxSyncTx) Commit(ctx context.Context) error {
	return database.WrapInSyncTxNoError(
		ctx,
		tx.syncTx,
		func(ctx context.Context) error {
			return tx.syncTx.Tx.Commit(ctx)
		},
	)
}

// Rollback is used to wrap [pgx.Tx.Rollback] in a [database.SyncTx].
func (tx *PgxSyncTx) Rollback(ctx context.Context) error {
	return database.WrapInSyncTxNoError(
		ctx,
		tx.syncTx,
		func(ctx context.Context) error {
			return tx.syncTx.Tx.Rollback(ctx)
		},
	)
}
