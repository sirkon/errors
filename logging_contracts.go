package errors

import (
	"go/token"
)

// ErrorContextConsumer is the contract that must be implemented by the logging side
// to receive the structured values stored in the error.
type ErrorContextConsumer interface {
	NextLink()
	Bool(name string, value bool)
	Int(name string, value int)
	Int8(name string, value int8)
	Int16(name string, value int16)
	Int32(name string, value int32)
	Int64(name string, value int64)
	Uint(name string, value uint)
	Uint8(name string, value uint8)
	Uint16(name string, value uint16)
	Uint32(name string, value uint32)
	Uint64(name string, value uint64)
	Flt32(name string, value float32)
	Flt64(name string, value float64)
	Str(name string, value string)
	Any(name string, value any)
	SetLinkInfo(loc token.Position, descr ErrorChainLinkDescriptor)
}

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

// ErrorChainLinkDescriptor describes a link in the processing chain. Implementations of this interface are used
// to describe the file:line locations in the error chain.
type ErrorChainLinkDescriptor interface {
	isErrorChainDescriptor()
}

// ErrorChainLinkWrap is given to chain links obtained via [Wrap]/[Wrapf].
type ErrorChainLinkWrap string

func (ErrorChainLinkWrap) isErrorChainDescriptor() {}

// ErrorChainLinkNew is given to links related to the creation of a new error via [New]/[Newf].
type ErrorChainLinkNew string

func (ErrorChainLinkNew) isErrorChainDescriptor() {}

// ErrorChainLinkContext is given to links where only structured context was added via [Just].
type ErrorChainLinkContext struct{}

func (ErrorChainLinkContext) isErrorChainDescriptor() {}
