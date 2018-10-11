package database

import (
	"errors"
	"fmt"
)

var (
	ErrSessionNotFound = errors.New("no session in database")
)

type UserNotFoundError struct {
	Field string
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("no user with this %v found", e.Field)
}
