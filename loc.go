package errors

import (
	"go/token"
	"runtime"
)

var insertLocations bool

// InsertLocations enables the insertion of error handling positions.
//
// WARNING!
//
// Calculating positions is a very expensive operation, so it is highly discouraged
// to enable position capturing in production. Enable only during local debugging.
//
// The recommended usage pattern is enabling when using "development" logging mode,
// with colored output and formatting intended for reading.
func InsertLocations() {
	insertLocations = true
}

// DoNotInsertLocations disables the insertion of error handling positions.
// This is the default mode.
func DoNotInsertLocations() {
	insertLocations = false
}

// setLoc adds the file:line location of error handling to the context.
//
//go:noinline
func (e *Error) setLoc() {
	_, fn, line, _ := runtime.Caller(2)
	e.loc = token.Position{Filename: fn, Line: line}
}
