package errors

import (
	"fmt"
)

// NewSentinel creates static error that can be a sentinel.
func NewSentinel(msg string) error {
	return &errorSentinel{
		value: msg,
	}
}

// NewSentinelf same as NewSentinel, just with a format instead of a static message.
func NewSentinelf(format string, a ...any) error {
	return &errorSentinel{
		value: fmt.Sprintf(format, a...),
	}
}

type errorSentinel struct {
	value string
}

func (e *errorSentinel) Error() string {
	return e.value
}
