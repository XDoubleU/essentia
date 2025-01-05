package postgres

import (
	"context"

	"github.com/XDoubleU/essentia/pkg/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// PgxSyncTx uses [database.SyncTx] to make sure
// [pgx.Tx] can be used concurrently.
type PgxSyncTx struct {
	syncTx *database.SyncTx[pgx.Tx]
}

// PgxSyncRow is a concurrent wrapper for [pgx.Row].
type PgxSyncRow struct {
	rows pgx.Rows
	err  error
}

// PgxSyncRows is a concurrent wrapper for [pgx.Rows].
type PgxSyncRows struct {
	values            [][]any
	rawValues         [][][]byte
	err               error
	fieldDescriptions []pgconn.FieldDescription
	commandTag        pgconn.CommandTag
	conn              *pgx.Conn
	i                 int
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
	return database.WrapInSyncTx(
		ctx,
		tx.syncTx,
		func(ctx context.Context) (*PgxSyncRows, error) {
			rows, err := tx.syncTx.Tx.Query(ctx, sql, args...)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			var results [][]any
			var rawResults [][][]byte
			for rows.Next() {
				var values []any
				values, err = rows.Values()
				if err != nil {
					break
				}

				temp := rows.RawValues()
				rawValues := make([][]byte, len(temp))
				copy(rawValues, temp)

				results = append(results, values)
				rawResults = append(rawResults, rawValues)
			}

			if err == nil {
				err = rows.Err()
			}

			return &PgxSyncRows{
				values:            results,
				rawValues:         rawResults,
				err:               err,
				fieldDescriptions: rows.FieldDescriptions(),
				commandTag:        rows.CommandTag(),
				conn:              rows.Conn(),
				i:                 -1,
			}, nil
		},
	)
}

// Close doesn't do anything for [PgxSyncRows] as these are closed in [Query].
func (rows *PgxSyncRows) Close() {
}

// CommandTag fetches the [pgconn.CommandTag].
func (rows *PgxSyncRows) CommandTag() pgconn.CommandTag {
	return rows.commandTag
}

// Conn fetches the [pgx.Conn].
func (rows *PgxSyncRows) Conn() *pgx.Conn {
	return rows.conn
}

// Err fetches any errors.
func (rows *PgxSyncRows) Err() error {
	return rows.err
}

// FieldDescriptions fetches [pgconn.FieldDescription]s.
func (rows *PgxSyncRows) FieldDescriptions() []pgconn.FieldDescription {
	return rows.fieldDescriptions
}

// Next continues to the next row of [PgxSyncRows] if there is one.
func (rows *PgxSyncRows) Next() bool {
	rows.i++
	return rows.i < len(rows.values)
}

// RawValues fetches the raw values of the current row.
func (rows *PgxSyncRows) RawValues() [][]byte {
	return rows.rawValues[rows.i]
}

// Scan scans the data of the current row into dest.
func (rows *PgxSyncRows) Scan(dest ...any) error {
	if err := rows.Err(); err != nil {
		return err
	}

	return pgx.ScanRow(
		pgtype.NewMap(),
		rows.FieldDescriptions(),
		rows.RawValues(),
		dest...)
}

// Values fetches the values of the current row.
func (rows *PgxSyncRows) Values() ([]any, error) {
	return rows.values[rows.i], nil
}

// QueryRow is used to wrap [pgx.Tx.QueryRow] in a [database.SyncTx].
func (tx *PgxSyncTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	rows, err := tx.Query(ctx, sql, args...)

	return &PgxSyncRow{
		rows: rows,
		err:  err,
	}
}

// Scan scans the data of [PgxSyncRow] into dest.
func (row *PgxSyncRow) Scan(dest ...any) error {
	if row.err != nil {
		return row.err
	}

	if err := row.rows.Err(); err != nil {
		return err
	}

	if !row.rows.Next() {
		return pgx.ErrNoRows
	}

	return row.rows.Scan(dest...)
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
