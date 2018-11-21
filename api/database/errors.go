package database

import (
	"errors"
	"fmt"
)

var (
	ErrNotNullConstraintViolation = errors.New("not null constraint violation")
	ErrUniqueConstraintViolation  = errors.New("unique constraint violation")
)

type UserNotFoundError struct {
	Field string
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("no user with this %v found", e.Field)
}
