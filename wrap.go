package errors

import (
	"fmt"
)

// Unwrap method needed for [errors.Is] and [errors.As] to work.
func (e *Error) Unwrap() error {
	return e.err
}

// Wrap annotates an error with a text message.
//
//go:noinline
func Wrap(err error, msg string) *Error {
	if err == nil {
		// TODO consider an option to create a fake error instead.
		_ = err.Error()
	}

	res := &Error{
		msg: msg,
		err: err,
		ctx: nil,
	}

	if insertLocations {
		res.setLoc()
	}

	return res
}

// Wrapf annotates an error with a formatted text message.
//
//go:noinline
func Wrapf(err error, format string, a ...any) *Error {
	if err == nil {
		// TODO consider an option to create a fake error instead.
		_ = err.Error()
	}

	res := &Error{
		msg: fmt.Sprintf(format, a...),
		err: err,
		ctx: nil,
	}

	if insertLocations {
		res.setLoc()
	}

	return res
}
