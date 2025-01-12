package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DB provides a uniform interface for the postgres
// database connection, pools and transactions.
type DB interface {
	Exec(
		ctx context.Context,
		sql string,
		arguments ...any,
	) (pgconn.CommandTag, error)
	Query(
		ctx context.Context,
		sql string,
		optionsAndArgs ...any,
	) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Ping(ctx context.Context) error
}
