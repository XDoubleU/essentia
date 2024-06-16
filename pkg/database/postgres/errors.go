package postgres

import (
	"errors"

	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func HandleError(err error) error {
	var pgxError *pgconn.PgError
	errors.As(err, &pgxError)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return httptools.ErrRecordNotFound
	case pgxError.Code == "23503":
		return httptools.ErrRecordNotFound
	case pgxError.Code == "23505":
		return httptools.ErrRecordUniqueValue
	default:
		return err
	}
}
