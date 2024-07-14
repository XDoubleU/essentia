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
		arguments ...interface{},
	) (pgconn.CommandTag, error)
	Query(
		ctx context.Context,
		sql string,
		optionsAndArgs ...interface{},
	) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
}
