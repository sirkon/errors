package errors

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

var _ error = Error{}

// Error implementation of error with structured context
type Error struct {
	msg string
	err error

	ctxPrefix string
	ctx       []contextTuple
}

func (e Error) Error() string {
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

type contextTuple struct {
	name  string
	value interface{}
}

// Is errors.Is support method
func (e Error) Is(err error) bool {
	if err == nil {
		return false
	}

	if v, ok := err.(Error); ok {
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

// As errors.As support method
func (e Error) As(target interface{}) bool {
	switch v := target.(type) {
	case *Error:
		*v = e
	case *errorContextDeliverer:
		*v = errorContextDeliverer{
			errCtx: e.ctx,
			next:   e.err,
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

// GetContextDeliverer extracts deliverer from err's chain if it does exist there
func GetContextDeliverer(err error) ErrorContextDeliverer {
	deliverer := Dig[errorContextDeliverer](err)
	if deliverer != nil {
		return *deliverer
	}

	return nil
}

type errorContextDeliverer struct {
	errCtx []contextTuple
	next   error
}

// Deliver to implement ErrorContextDeliverer
func (e errorContextDeliverer) Deliver(cons ErrorContextConsumer) {
	for _, item := range e.errCtx {
		switch v := item.value.(type) {
		case bool:
			cons.Bool(item.name, v)
		case int:
			cons.Int(item.name, v)
		case int8:
			cons.Int8(item.name, v)
		case int16:
			cons.Int16(item.name, v)
		case int32:
			cons.Int32(item.name, v)
		case int64:
			cons.Int64(item.name, v)
		case uint:
			cons.Uint(item.name, v)
		case uint8:
			cons.Uint8(item.name, v)
		case uint16:
			cons.Uint16(item.name, v)
		case uint32:
			cons.Uint32(item.name, v)
		case uint64:
			cons.Uint64(item.name, v)
		case float32:
			cons.Float32(item.name, v)
		case float64:
			cons.Float64(item.name, v)
		case string:
			cons.String(item.name, v)
		case []byte:
			if utf8.Valid(v) {
				cons.String(item.name, string(v))
				break
			}
			cons.String(item.name, base64.StdEncoding.EncodeToString(v))
		case fmt.Stringer:
			cons.String(item.name, v.String())
		default:
			cons.Any(item.name, item.value)
		}
	}

	if v := Dig[errorContextDeliverer](e.next); v != nil {
		(*v).Deliver(cons)
	}
}

func (e errorContextDeliverer) Error() string {
	return ""
}
