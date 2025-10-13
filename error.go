package errors

import (
	"errors"
	"go/token"
	"math"
	"strings"
)

var _ error = new(Error)

// Error is an error implementation with structured context support.
type Error struct {
	msg string
	err error
	loc token.Position

	ctxPrefix string
	ctx       []contextTuple
}

func (e *Error) Error() string {
	var buf strings.Builder
	if e.msg != "" {
		buf.WriteString(e.msg)
		if e.err == nil {
			return buf.String()
		}
		buf.WriteString(": ")
	}

	buf.WriteString(e.err.Error())
	return buf.String()
}

// Is implements support for [Is].
func (e *Error) Is(err error) bool {
	if err == nil {
		return false
	}

	if v, ok := err.(*Error); ok {
		if e.err == nil {
			return e.msg == v.msg && v.err == nil
		}

		if e.msg == v.msg && errors.Is(e.err, v.err) {
			return true
		}
	}

	if e.err == nil {
		return false
	}

	return errors.Is(e.err, err)
}

// As implements support for [As].
func (e *Error) As(target any) bool {
	switch v := target.(type) {
	case **Error:
		*v = e
	case *errorContextDeliverer:
		*v = errorContextDeliverer{
			text:   e.msg,
			errCtx: e.ctx,
			next:   e.err,
			loc:    e.loc,
		}

		return true
	default:
		if e.err == nil {
			return false
		}

		return errors.As(e.err, target)
	}

	return false
}

// GetContextDeliverer returns the structured-context deliverer for the given error.
func GetContextDeliverer(err error) ErrorContextDeliverer {
	deliverer, ok := AsType[errorContextDeliverer](err)
	if !ok {
		return nil
	}

	return deliverer
}

type errorContextDeliverer struct {
	text   string
	errCtx []contextTuple
	loc    token.Position
	next   error
}

// Deliver implements ErrorContextDeliverer.
func (e errorContextDeliverer) Deliver(cons ErrorContextConsumer) {
	cons.NextLink()
	for _, item := range e.errCtx {
		switch item.kind {
		case tupleKindBool:
			var v bool
			if item.scalar > 0 {
				v = true
			}
			cons.Bool(item.name, v)
		case tupleKindInt:
			cons.Int(item.name, int(item.scalar))
		case tupleKindInt8:
			cons.Int8(item.name, int8(item.scalar))
		case tupleKindInt16:
			cons.Int16(item.name, int16(item.scalar))
		case tupleKindInt32:
			cons.Int32(item.name, int32(item.scalar))
		case tupleKindInt64:
			cons.Int64(item.name, int64(item.scalar))
		case tupleKindUint:
			cons.Uint(item.name, uint(item.scalar))
		case tupleKindUint8:
			cons.Uint8(item.name, uint8(item.scalar))
		case tupleKindUint16:
			cons.Uint16(item.name, uint16(item.scalar))
		case tupleKindUint32:
			cons.Uint32(item.name, uint32(item.scalar))
		case tupleKindUint64:
			cons.Uint64(item.name, item.scalar)
		case tupleKindFloat32:
			cons.Flt32(item.name, math.Float32frombits(uint32(item.scalar)))
		case tupleKindFloat64:
			cons.Flt64(item.name, math.Float64frombits(item.scalar))
		case tupleKindString:
			cons.Str(item.name, item.str)
		case tupleKindStrings, tupleKindAny:
			cons.Any(item.name, item.value)
		case tupleKindInvalid:
			cons.Any(item.name, item.value)
		}
	}

	// Add the location info if it is set.
	var descr ErrorChainLinkDescriptor
	switch {
	case e.next == nil:
		descr = ErrorChainLinkNew(e.text)
	case e.text == "":
		descr = ErrorChainLinkContext{}
	default:
		descr = ErrorChainLinkWrap(e.text)
	}
	cons.SetLinkInfo(e.loc, descr)

	if v, ok := AsType[errorContextDeliverer](e.next); ok {
		v.Deliver(cons)
	}
}

func (e errorContextDeliverer) Error() string {
	return ""
}
