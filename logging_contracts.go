package errors

// ErrorContextDeliverer is what is returned by the [GetContextDeliverer] function.
// Example of how to work with slog.
//
//	type slogConsumer struct{
//	   fields []any
//	}
//
//	func (c *slogConsumer) Int(name, value int) {
//	   c.fields = append(c.fields, slog.Int(name, value))
//	}
//
//	// And so on for the rest of the methods to implement errors.ErrorContextConsumer.
//
//	// Log logs at the Info level.
//	func Log(msg string, fields []slog.Attr) {
//	    var attrs []any
//	    for _, field := range fields {
//	        // Save all original fields.
//	        attrs = append(attrs, field)
//
//	        // Need to unwrap errors. For this, we look for fields that contain them.
//	        fv := field.Value.Any()
//	        e, ok := fv.(error)
//	        if !ok {
//	            continue
//	        }
//
//	        // v contains an error. Try to get context from it.
//	        dlr := errors.GetContextDeliverer(e)
//	        if dlr == nil {
//	            // This is not our error.
//	            continue
//	        }
//
//	        // Get the context and add the extracted fields.
//	        var errCtx slogConsumer{}
//	        dlr.Deliver(&errCtx)
//	        attrs = append(attrs, errCtx.fields...)
//	    }
//
//	    slog.Info(msg, attrs...)
//	}
//
// A working example can be found in ./internal/example/main.go.
type ErrorContextDeliverer interface {
	Deliver(cons ErrorContextConsumer)
	Error() string
}

// ErrorContextConsumer is the contract that must be implemented by the logging side
// to receive the structured values stored in the error.
type ErrorContextConsumer interface {
	New(msg string) ErrorContextBuilder
	Wrap(msg string) ErrorContextBuilder
	Just() ErrorContextBuilder
}

// ErrorContextBuilder is the contract for dispatching context values
// per layer of processing.
type ErrorContextBuilder interface {
	Bool(name string, value bool)
	Int64(name string, value int64)
	Uint64(name string, value uint64)
	Flt64(name string, value float64)
	Str(name string, value string)
	Any(name string, value any)
	Loc(position string)
	Finalize()
}

// GetContextDeliverer returns the structured-context deliverer for the given error.
func GetContextDeliverer(err error) ErrorContextDeliverer {
	deliverer, ok := AsType[*errorContextDeliverer](err)
	if !ok {
		return nil
	}

	return deliverer
}

// MustGetContextDeliverer guaranteed delivery on the instance of Error istelf.
func MustGetContextDeliverer(err *Error) ErrorContextDeliverer {
	return &errorContextDeliverer{
		tgt: err,
	}
}
