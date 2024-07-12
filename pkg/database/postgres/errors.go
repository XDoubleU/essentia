package postgres

import (
	"errors"

	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func PgxErrorToHTTPError(err error) error {
	var pgxError *pgconn.PgError
	errors.As(err, &pgxError)

	switch {
	case errors.Is(err, pgx.ErrNoRows), pgxError.Code == pgerrcode.ForeignKeyViolation:
		return httptools.ErrRecordNotFound
	case pgxError.Code == pgerrcode.UniqueViolation:
		return httptools.ErrRecordUniqueValue
	default:
		return err
	}
}
