package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xdoubleu/essentia/pkg/database"
)

type SpanDB struct {
	DB     DB
	dbName string
}

func NewSpanDB(db DB) SpanDB {
	return SpanDB{
		DB:     db,
		dbName: "postgresql",
	}
}

func (db SpanDB) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (pgconn.CommandTag, error) {
	return database.WrapWithSpan(ctx, db.dbName, db.DB.Exec, sql, arguments...)
}

func (db SpanDB) Query(
	ctx context.Context,
	sql string,
	optionsAndArgs ...any,
) (pgx.Rows, error) {
	return database.WrapWithSpan(ctx, db.dbName, db.DB.Query, sql, optionsAndArgs...)
}

func (db SpanDB) QueryRow(
	ctx context.Context,
	sql string,
	optionsAndArgs ...any,
) pgx.Row {
	return database.WrapWithSpanNoError(ctx, db.dbName, db.DB.QueryRow, sql, optionsAndArgs...)
}

func (db SpanDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.DB.Begin(ctx)
}
