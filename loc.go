package errors

import (
	"log/slog"
	"runtime"
	"strconv"
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
func (e *Error) setLoc(skip int) {
	_, fn, line, _ := runtime.Caller(skip)
	e.attrs = append(e.attrs, errorAttr{
		kind:  errorAttrKindLoc,
		value: slog.StringValue(fn + ":" + strconv.Itoa(line)),
	})
}
