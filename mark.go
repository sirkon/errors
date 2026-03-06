package errors

import (
	"log/slog"
)

// Mark marks an error with given "mark" which is not shown in an output.
func Mark(err error, mark any) *Error {
	if e, ok := err.(*Error); ok {
		e.attrs = append(e.attrs, errorAttr{
			value: slog.AnyValue(mark),
			kind:  errorAttrKindMarker,
		})
		return e
	}

	attrs := make([]errorAttr, 0, errorContextLengthPrediction)
	attrs = append(
		attrs,
		errorAttr{
			key:   "",
			value: slog.AnyValue(err),
			kind:  errorAttrKindPhantomJust,
		},
		errorAttr{
			value: slog.AnyValue(mark),
			kind:  errorAttrKindMarker,
		},
	)
	return &Error{
		attrs: attrs,
	}
}

func HasMark[T any](err error) (v T, ok bool) {
	e, ok := err.(*Error)
	if !ok {
		e, ok = AsType[*Error](err)
		if !ok {
			var zero T
			return zero, false
		}
	}

	var wrappedErr error
	for _, attr := range e.attrs {
		switch attr.kind {
		case errorAttrKindMarker:
			v, ok := attr.value.Any().(T)
			if ok {
				return v, true
			}
		case errorAttrKindOutterWrap, errorAttrKindOutterJust:
			wrappedErr = attr.value.Any().(error)
		}
	}

	if wrappedErr == nil {
		var zero T
		return zero, false
	}

	return HasMark[T](wrappedErr)
}
