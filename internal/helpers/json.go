package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

func ReadJSON[T any](body io.Reader, dst T, allowEmpty bool) error {
	err := json.NewDecoder(body).Decode(dst)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf(
				"body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf(
					"body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field,
				)
			}

			return fmt.Errorf(
				"body contains incorrect JSON type (at character %d)",
				unmarshalTypeError.Offset,
			)

		case errors.Is(err, io.EOF):
			if allowEmpty {
				return nil
			}

			return errors.New("body must not be empty")

		case errors.As(err, &invalidUnmarshalError):
			return err

		default:
			return err
		}
	}

	return nil
}
