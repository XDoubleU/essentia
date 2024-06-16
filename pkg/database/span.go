package database

import (
	"context"

	"github.com/getsentry/sentry-go"
)

func StartSpan(ctx context.Context, databaseName string, sql string) *sentry.Span {
	span := sentry.StartSpan(ctx, "db.query", sentry.WithDescription(sql))
	span.SetData("db.system", databaseName)

	return span
}
