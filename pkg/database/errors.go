package database

import "errors"

var (
	// ErrResourceNotFound is an error with value "resource not found".
	ErrResourceNotFound = errors.New("resource not found")
	// ErrResourceConflict is an error with value "resource unique value already used".
	ErrResourceConflict = errors.New("resource conflicts with existing resource")
)
