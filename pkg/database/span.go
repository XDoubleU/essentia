package database

import (
	"context"

	"github.com/getsentry/sentry-go"
)

// StartSpan is used to start a [sentry.Span].
func StartSpan(ctx context.Context, dbName string, sql string) *sentry.Span {
	span := sentry.StartSpan(ctx, "db.query", sentry.WithDescription(sql))
	span.SetData("db.system", dbName)

	return span
}

// WrapWithSpan is used to wrap a
// database action in a [sentry.Span].
func WrapWithSpan[T any](
	ctx context.Context,
	dbName string,
	queryFunc func(ctx context.Context, sql string, args ...any) (T, error),
	sql string, args ...any) (T, error) {
	span := StartSpan(ctx, dbName, sql)
	defer span.Finish()

	return queryFunc(ctx, sql, args...)
}

// WrapWithSpanNoError is used to wrap a
// database action in a [sentry.Span].
// The executed database action shouldn't return an error.
func WrapWithSpanNoError[T any](
	ctx context.Context,
	dbName string,
	queryFunc func(ctx context.Context, sql string, args ...any) T,
	sql string, args ...any) T {
	span := StartSpan(ctx, dbName, sql)
	defer span.Finish()

	return queryFunc(ctx, sql, args...)
}
