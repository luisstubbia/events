package model

import (
	"errors"
	"fmt"
	"time"
)

type ValidationError struct {
	Message string
	When    time.Time
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s (at: %s)", e.Message, e.When.Format("2006-01-02 15:04:05"))
}

func (e *ValidationError) Is(err error) bool {
	var validationError *ValidationError
	ok := errors.As(err, &validationError)
	return ok
}

func NewValidationError(message string, when time.Time) error {
	return &ValidationError{
		Message: message,
		When:    when,
	}
}

func NewValidationErrorWithTime(message string) error {
	return &ValidationError{
		Message: message,
		When:    time.Now().UTC(),
	}
}
