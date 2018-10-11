package handlers

import (
	"fmt"
)

type ParseJSONError struct {
	msg error
}

func (e ParseJSONError) Error() string {
	return fmt.Sprintf("error while parsing JSON: %v", e.msg)
}
