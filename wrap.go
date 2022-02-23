package errors

import "fmt"

// Unwrap returns naked error out of these wraps
func (e Error) Unwrap() error {
	return e.err
}

// Wrap constructs a new error by wrapping given message over the existing one
func Wrap(err error, msg string) Error {
	if err == nil {
		// this is intentional, you must not wrap nil error
		err.Error()
	}

	return Error{
		msg: msg,
		err: err,
		ctx: nil,
	}
}

// Wrapf calls Wrap function with a message built using given format
func Wrapf(err error, format string, a ...interface{}) Error {
	return Wrap(err, fmt.Sprintf(format, a...))
}
