package errors

import (
	"fmt"
)

// New creates new error with the given text.
func New(msg string) *Error {
	res := &Error{
		attrs: make([]errorAttr, 0, errorContextLengthPrediction),
	}
	res.attrs = append(res.attrs, errorAttr{
		kind: errorAttrKindNew,
		key:  msg,
	})

	if insertLocations {
		res.setLoc(2)
	}

	return res
}

// Newf creates new error with the given format.
func Newf(format string, a ...any) *Error {
	res := &Error{
		attrs: make([]errorAttr, 0, errorContextLengthPrediction),
	}
	res.attrs = append(res.attrs, errorAttr{
		kind: errorAttrKindNew,
		key:  fmt.Sprintf(format, a...),
	})

	if insertLocations {
		res.setLoc(2)
	}

	return res
}
