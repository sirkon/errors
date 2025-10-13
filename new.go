package errors

import (
	"fmt"
)

// New creates a new error with the given message.
//
//go:noinline
func New(msg string) *Error {
	res := &Error{
		msg: msg,
	}

	if insertLocations {
		res.setLoc()
	}

	return res
}

// Newf creates a new error with the given formatted message.
//
//go:noinline
func Newf(format string, a ...any) *Error {
	res := &Error{
		msg: fmt.Sprintf(format, a...),
	}

	if insertLocations {
		res.setLoc()
	}

	return res
}
