package errors

import "fmt"

// New creates new Error with a given message
func New(msg string) Error {
	return Error{
		msg: msg,
	}
}

// Newf same as New, with formatted error message
func Newf(format string, a ...interface{}) Error {
	return Error{
		msg: fmt.Sprintf(format, a...),
	}
}
