package database

import (
	"context"

	"github.com/getsentry/sentry-go"
)

func startSpan(ctx context.Context, dbName string, sql string) *sentry.Span {
	span := sentry.StartSpan(ctx, "db.query", sentry.WithDescription(sql))
	span.SetData("db.system", dbName)

	return span
}

func WrapWithSpan[T any](
	ctx context.Context,
	dbName string,
	queryFunc func(ctx context.Context, sql string, args ...any) (T, error),
	sql string, args ...any) (T, error) {
	span := startSpan(ctx, dbName, sql)
	defer span.Finish()

	return queryFunc(ctx, sql, args...)
}

func WrapWithSpanNoError[T any](
	ctx context.Context,
	dbName string,
	queryFunc func(ctx context.Context, sql string, args ...any) T,
	sql string, args ...any) T {
	span := startSpan(ctx, dbName, sql)
	defer span.Finish()

	return queryFunc(ctx, sql, args...)
}
