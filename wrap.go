package errors

import (
	"fmt"
	"log/slog"
)

// Wrap annotates given error with the given message.
func Wrap(err error, msg string) *Error {
	return wrap(err, msg)
}

// Wrapf annotates given error with the given format.
func Wrapf(err error, format string, a ...any) *Error {
	return wrap(err, fmt.Sprintf(format, a...))
}

// Just возвращает *Error позволяющий добавить контекст данной ошибке.
func Just(err error) *Error {
	e, ok := err.(*Error)
	if !ok {
		res := &Error{
			attrs: make([]errorAttr, 0, errorContextLengthPrediction),
		}
		res.attrs = append(res.attrs, errorAttr{
			kind:  errorAttrKindOutterJust,
			value: slog.AnyValue(err),
		})
		return res
	}

	e.attrs = append(e.attrs, errorAttr{
		kind: errorAttrKindJust,
	})

	if insertLocations {
		e.setLoc(2)
	}

	return e
}

func wrap(err error, msg string) *Error {
	e, ok := err.(*Error)
	if !ok {
		res := &Error{
			attrs: make([]errorAttr, 0, errorContextLengthPrediction),
		}
		res.attrs = append(res.attrs, errorAttr{
			kind:  errorAttrKindOutterWrap,
			key:   msg,
			value: slog.AnyValue(err),
		})
		return res
	}

	e.attrs = append(e.attrs, errorAttr{
		kind: errorAttrKindWrap,
		key:  msg,
	})

	if insertLocations {
		e.setLoc(3)
	}

	return e
}
