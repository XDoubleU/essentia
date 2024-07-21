package postgres

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xdoubleu/essentia/pkg/httptools"
)

// PgxErrorToHTTPError converts a database error
// from [github.com/jackc/pgx] to an appropriate HTTP error.
func PgxErrorToHTTPError(err error) error {
	var pgxError *pgconn.PgError
	errors.As(err, &pgxError)

	switch {
	case errors.Is(err, pgx.ErrNoRows), pgxError.Code == pgerrcode.ForeignKeyViolation:
		return httptools.ErrResourceNotFound
	case pgxError.Code == pgerrcode.UniqueViolation:
		return httptools.ErrResourceUniqueValue
	default:
		return err
	}
}
