package errors

import (
	"errors"
	"fmt"

	"github.com/xdoubleu/essentia/internal/shared"
)

var (
	ErrFailedValidation = errors.New("failed validation")
	ErrUnauthorized     = errors.New("unauthorized")
)

type NotFoundError struct {
	resourceName    string
	identifierValue string
	JsonField       string
}

type ConflictError struct {
	resourceName    string
	identifierValue string
	JsonField       string
}

type BadRequestError struct {
	err error
}

func NewNotFoundError(resourceName string, identifierValue any, jsonField string) NotFoundError {
	value, err := shared.AnyToString(identifierValue)
	if err != nil {
		panic(err)
	}

	return NotFoundError{
		resourceName:    resourceName,
		identifierValue: value,
		JsonField:       jsonField,
	}
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf(
		"%s with %s '%s' doesn't exist",
		err.resourceName,
		err.JsonField,
		err.identifierValue,
	)
}

func NewConflictError(resourceName string, identifierValue any, jsonField string) ConflictError {
	value, err := shared.AnyToString(identifierValue)
	if err != nil {
		panic(err)
	}

	return ConflictError{
		resourceName:    resourceName,
		identifierValue: value,
		JsonField:       jsonField,
	}
}

func (err ConflictError) Error() string {
	return fmt.Sprintf(
		"%s with %s '%s' already exists",
		err.resourceName,
		err.JsonField,
		err.identifierValue,
	)
}

func NewBadRequestError(err error) BadRequestError {
	return BadRequestError{
		err: err,
	}
}

func (err BadRequestError) Error() string {
	return err.Error()
}
