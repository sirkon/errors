package errors

import (
	"fmt"
	"math"
	"strings"
	"unicode"
	"unicode/utf8"
)

type tupleKind uint64

const (
	tupleKindInvalid tupleKind = iota
	tupleKindBool
	tupleKindInt
	tupleKindInt8
	tupleKindInt16
	tupleKindInt32
	tupleKindInt64
	tupleKindUint
	tupleKindUint8
	tupleKindUint16
	tupleKindUint32
	tupleKindUint64
	tupleKindFloat32
	tupleKindFloat64
	tupleKindString
	tupleKindStrings
	tupleKindAny
)

// We adopt the approach used by structured loggers to help reduce the number of allocations.
// This is not critical right now, but it may become relevant if we optimize for fewer allocations later.
// In that case, this approach will fit naturally.
type contextTuple struct {
	name   string
	kind   tupleKind
	scalar uint64
	str    string
	value  any
}

// Pfx sets (or replaces) the prefix that will be prepended to subsequent context field names.
func (e *Error) Pfx(prefix string) *Error {
	if prefix == "" {
		return e
	}

	e.ctxPrefix = strings.TrimSuffix(prefix, "-") + "-"
	return e
}

// WithCtx merges the given [Context] into the error.
// The current error's prefix is prepended to the names of the added fields.
func (e *Error) WithCtx(ctx *Context) *Error {
	fields := e.ctx
	for _, field := range ctx.fields {
		fields = append(fields, contextTuple{
			name:   e.ctxPrefix + field.name,
			kind:   field.kind,
			scalar: field.scalar,
			str:    field.str,
			value:  field.value,
		})
	}

	return &Error{
		msg:       e.msg,
		err:       e.err,
		loc:       e.loc,
		ctxPrefix: e.ctxPrefix,
		ctx:       fields,
	}
}

// Bool adds a named boolean value to the context.
func (e *Error) Bool(name string, value bool) *Error {
	var v uint64
	if value {
		v = 1
	}
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:   e.ctxPrefix + name,
			kind:   tupleKindBool,
			scalar: v,
		}),
		loc: e.loc,
	}
}

// Int adds a named integer value (int) to the context.
func (e *Error) Int(name string, value int) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:   e.ctxPrefix + name,
			kind:   tupleKindInt,
			scalar: uint64(value),
		}),
		loc: e.loc,
	}
}

// Int8 adds a named integer value (int8) to the context.
func (e *Error) Int8(name string, value int8) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:   e.ctxPrefix + name,
			kind:   tupleKindInt8,
			scalar: uint64(value),
		}),
		loc: e.loc,
	}
}

// Int16 adds a named integer value (int16) to the context.
func (e *Error) Int16(name string, value int16) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:   e.ctxPrefix + name,
			kind:   tupleKindInt16,
			scalar: uint64(value),
		}),
		loc: e.loc,
	}
}

// Int32 adds a named integer value (int32) to the context.
func (e *Error) Int32(name string, value int32) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:   e.ctxPrefix + name,
			kind:   tupleKindInt32,
			scalar: uint64(value),
		}),
		loc: e.loc,
	}
}

// Int64 adds a named integer value (int64) to the context.
func (e *Error) Int64(name string, value int64) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:   e.ctxPrefix + name,
			kind:   tupleKindInt64,
			scalar: uint64(value),
		}),
		loc: e.loc,
	}
}

// Uint adds a named unsigned integer value (uint) to the context.
func (e *Error) Uint(name string, value uint) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:   e.ctxPrefix + name,
			kind:   tupleKindUint,
			scalar: uint64(value),
		}),
		loc: e.loc,
	}
}

// Uint8 adds a named unsigned integer value (uint8) to the context.
func (e *Error) Uint8(name string, value uint8) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:   e.ctxPrefix + name,
			kind:   tupleKindUint8,
			scalar: uint64(value),
		}),
		loc: e.loc,
	}
}

// Uint16 adds a named unsigned integer value (uint16) to the context.
func (e *Error) Uint16(name string, value uint16) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:   e.ctxPrefix + name,
			kind:   tupleKindUint16,
			scalar: uint64(value),
		}),
		loc: e.loc,
	}
}

// Uint32 adds a named unsigned integer value (uint32) to the context.
func (e *Error) Uint32(name string, value uint32) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:   e.ctxPrefix + name,
			kind:   tupleKindUint32,
			scalar: uint64(value),
		}),
		loc: e.loc,
	}
}

// Uint64 adds a named unsigned integer value (uint64) to the context.
func (e *Error) Uint64(name string, value uint64) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:   e.ctxPrefix + name,
			kind:   tupleKindUint64,
			scalar: value,
		}),
		loc: e.loc,
	}
}

// Flt32 adds a named floating-point value (float32) to the context.
func (e *Error) Flt32(name string, value float32) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:   e.ctxPrefix + name,
			kind:   tupleKindFloat32,
			scalar: uint64(math.Float32bits(value)),
		}),
		loc: e.loc,
	}
}

// Flt64 adds a named floating-point value (float64) to the context.
func (e *Error) Flt64(name string, value float64) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:   e.ctxPrefix + name,
			kind:   tupleKindFloat64,
			scalar: math.Float64bits(value),
		}),
		loc: e.loc,
	}
}

// Str adds a named string value to the context.
func (e *Error) Str(name string, value string) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name: e.ctxPrefix + name,
			kind: tupleKindString,
			str:  value,
		}),
		loc: e.loc,
	}
}

// Stg adds a named value implementing fmt.Stringer to the context.
func (e *Error) Stg(name string, value fmt.Stringer) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name: e.ctxPrefix + name,
			kind: tupleKindString,
			str:  value.String(),
		}),
		loc: e.loc,
	}
}

// Strs adds a named slice of strings value to the context.
func (e *Error) Strs(name string, value []string) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			kind:  tupleKindStrings,
			value: value,
		}),
		loc: e.loc,
	}
}

// Bytes adds a named slice of bytes value to the context.
//
// Attention:
//
// This operation can be computationally heavy,
// as it evaluates whether the sequence can be represented
// as a string, after which it is saved in either string or object format.
func (e *Error) Bytes(name string, value []byte) *Error {
	tuple := contextTuple{
		name: e.ctxPrefix + name,
	}
	if isPrintableStringWithSpaces(value) {
		tuple.kind = tupleKindString
		tuple.str = string(value)
	} else {
		tuple.kind = tupleKindAny
		tuple.value = value
	}

	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx:       append(e.ctx, tuple),
		loc:       e.loc,
	}
}

// Type adds a type name to the context.
func (e *Error) Type(name string, typ any) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name: e.ctxPrefix + name,
			kind: tupleKindString,
			str:  fmt.Sprintf("%T", typ),
		}),
		loc: e.loc,
	}
}

// Any adds an arbitrary named value to the context.
func (e *Error) Any(name string, value any) *Error {
	return &Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			kind:  tupleKindAny,
			value: value,
		}),
		loc: e.loc,
	}
}

func isPrintableStringWithSpaces(b []byte) bool {
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		if r == utf8.RuneError || (!unicode.IsPrint(r) && !unicode.IsSpace(r)) {
			return false
		}
		b = b[size:]
	}
	return true
}
