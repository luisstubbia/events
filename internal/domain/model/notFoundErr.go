package model

import (
	"errors"
	"fmt"
	"time"
)

type NotFoundError struct {
	Message string
	When    time.Time
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s (at: %s)", e.Message, e.When.Format("2006-01-02 15:04:05"))
}

func (e *NotFoundError) Is(err error) bool {
	var notFoundErr *NotFoundError
	ok := errors.As(err, &notFoundErr)
	return ok
}

func NewNotFoundError(message string, when time.Time) error {
	return &NotFoundError{
		Message: message,
		When:    when,
	}
}
