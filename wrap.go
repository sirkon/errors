package errors

import (
	"fmt"
	"strings"
)

var _ error = &wrappedError{}

type wrappedError struct {
	msgs []string
	err  error
}

func (w *wrappedError) Error() string {
	var buf strings.Builder
	for i := len(w.msgs) - 1; i >= 0; i-- {
		buf.WriteString(w.msgs[i])
		buf.WriteString(": ")
	}
	buf.WriteString(w.err.Error())
	return buf.String()
}

// Unwrap returns naked error out of these wraps
func (w *wrappedError) Unwrap() error {
	return w.err
}

// Wrap consturcts a new error by wrapping given message into an error
func Wrap(err error, msg string) error {
	if err == nil {
		// this is intentional
		err.Error()
	}
	switch v := err.(type) {
	case *wrappedError:
		v.msgs = append(v.msgs, msg)
		return v
	default:
		msgs := make([]string, 1, 4)
		msgs[0] = msg
		return &wrappedError{
			msgs: msgs,
			err:  err,
		}
	}
}

// Wrapf calls Wrap function with a message built with given format
func Wrapf(err error, format string, a ...interface{}) error {
	return Wrap(err, fmt.Sprintf(format, a...))
}
