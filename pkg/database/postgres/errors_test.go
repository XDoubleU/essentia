package postgres_test

import (
	"testing"

	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func newPgError(code string) *pgconn.PgError {
	//nolint:exhaustruct //not using other fields
	return &pgconn.PgError{
		Code: code,
	}
}

func TestErrResourceNotFound(t *testing.T) {
	err1 := postgres.PgxErrorToHTTPError(pgx.ErrNoRows)
	err2 := postgres.PgxErrorToHTTPError(
		newPgError(pgerrcode.ForeignKeyViolation),
	)

	assert.ErrorIs(t, err1, httptools.ErrResourceNotFound)
	assert.ErrorIs(t, err2, httptools.ErrResourceNotFound)
}

func TestErrResourceUniqueValue(t *testing.T) {
	err := postgres.PgxErrorToHTTPError(
		newPgError(pgerrcode.UniqueViolation),
	)

	assert.ErrorIs(t, err, httptools.ErrResourceUniqueValue)
}
