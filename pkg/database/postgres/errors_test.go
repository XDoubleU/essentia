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

func TestErrRecordNotFound(t *testing.T) {
	err1 := postgres.PgxErrorToHTTPError(pgx.ErrNoRows)
	err2 := postgres.PgxErrorToHTTPError(
		&pgconn.PgError{Code: pgerrcode.ForeignKeyViolation},
	)

	assert.ErrorIs(t, err1, httptools.ErrRecordNotFound)
	assert.ErrorIs(t, err2, httptools.ErrRecordNotFound)
}

func TestErrRecordUniqueValue(t *testing.T) {
	err := postgres.PgxErrorToHTTPError(
		&pgconn.PgError{Code: pgerrcode.UniqueViolation},
	)

	assert.ErrorIs(t, err, httptools.ErrRecordUniqueValue)
}
