package errors

import (
	"fmt"
	"math"
)

// Ctx creates an empty Context.
func Ctx() *Context {
	return &Context{}
}

// CtxFrom creates a new context by copying values from an existing one.
func CtxFrom(ctx *Context) *Context {
	return Ctx().WithCtx(ctx)
}

// Context is an entity for implementing context composition.
// Using this object implies mutability; for a reuse approach, use [Context.WithCtx].
type Context struct {
	fields []contextTuple
}

// WithCtx copies values from the given context into the current one.
//
// This method supports logic where multiple contexts can be used,
// which may have common parts but differ in others.
// Something like:
//
//	ctx1 := errors.Ctx().Int("code", 404)
//	ctx2 := ctx1.Str("msg", "not found")
//	ctx3 := ctx1.Bool("retried", isRetried)
//
// Will not work. Here ctx1, ctx2, and ctx3 will be pointers to the same
// object and, accordingly, will yield the same context.
//
// See the example for detailed usage.
func (b *Context) WithCtx(ctx *Context) *Context {
	b.fields = append(b.fields, ctx.fields...)
	return b
}

// Bool adds a named boolean value (bool) to the context builder.
func (b *Context) Bool(name string, value bool) *Context {
	var val uint64
	if value {
		val = 1
	}
	b.fields = append(b.fields, contextTuple{
		name:   name,
		kind:   tupleKindBool,
		scalar: val,
		str:    "",
		value:  nil,
	})

	return b
}

// Int adds a named integer value (int) to the context builder.
func (b *Context) Int(name string, value int) *Context {
	b.fields = append(b.fields, contextTuple{
		name:   name,
		kind:   tupleKindInt,
		scalar: uint64(value),
	})
	return b
}

// Int8 adds a named integer value (int8) to the context builder.
func (b *Context) Int8(name string, value int8) *Context {
	b.fields = append(b.fields, contextTuple{
		name:   name,
		kind:   tupleKindInt8,
		scalar: uint64(value),
	})
	return b
}

// Int16 adds a named integer value (int16) to the context builder.
func (b *Context) Int16(name string, value int16) *Context {
	b.fields = append(b.fields, contextTuple{
		name:   name,
		kind:   tupleKindInt16,
		scalar: uint64(value),
	})
	return b
}

// Int32 adds a named integer value (int32) to the context builder.
func (b *Context) Int32(name string, value int32) *Context {
	b.fields = append(b.fields, contextTuple{
		name:   name,
		kind:   tupleKindInt32,
		scalar: uint64(value),
	})
	return b
}

// Int64 adds a named integer value (int64) to the context builder.
func (b *Context) Int64(name string, value int64) *Context {
	b.fields = append(b.fields, contextTuple{
		name:   name,
		kind:   tupleKindInt64,
		scalar: uint64(value),
	})
	return b
}

// Uint adds a named unsigned integer value (uint) to the context builder.
func (b *Context) Uint(name string, value uint) *Context {
	b.fields = append(b.fields, contextTuple{
		name:   name,
		kind:   tupleKindUint,
		scalar: uint64(value),
	})
	return b
}

// Uint8 adds a named unsigned integer value (uint8) to the context builder.
func (b *Context) Uint8(name string, value uint8) *Context {
	b.fields = append(b.fields, contextTuple{
		name:   name,
		kind:   tupleKindUint8,
		scalar: uint64(value),
	})
	return b
}

// Uint16 adds a named unsigned integer value (uint16) to the context builder.
func (b *Context) Uint16(name string, value uint16) *Context {
	b.fields = append(b.fields, contextTuple{
		name:   name,
		kind:   tupleKindUint16,
		scalar: uint64(value),
	})
	return b
}

// Uint32 adds a named unsigned integer value (uint32) to the context builder.
func (b *Context) Uint32(name string, value uint32) *Context {
	b.fields = append(b.fields, contextTuple{
		name:   name,
		kind:   tupleKindUint32,
		scalar: uint64(value),
	})
	return b
}

// Uint64 adds a named unsigned integer value (uint64) to the context builder.
func (b *Context) Uint64(name string, value uint64) *Context {
	b.fields = append(b.fields, contextTuple{
		name:   name,
		kind:   tupleKindUint64,
		scalar: value,
	})
	return b
}

// Flt32 adds a named floating-point value (float32) to the context builder.
func (b *Context) Flt32(name string, value float32) *Context {
	b.fields = append(b.fields, contextTuple{
		name:   name,
		kind:   tupleKindFloat32,
		scalar: uint64(math.Float32bits(value)),
	})
	return b
}

// Flt64 adds a named floating-point value (float64) to the context builder.
func (b *Context) Flt64(name string, value float64) *Context {
	b.fields = append(b.fields, contextTuple{
		name:   name,
		kind:   tupleKindFloat64,
		scalar: math.Float64bits(value),
	})
	return b
}

// Str adds a named string value to the context builder.
func (b *Context) Str(name string, value string) *Context {
	b.fields = append(b.fields, contextTuple{
		name: name,
		kind: tupleKindString,
		str:  value,
	})
	return b
}

// Stg adds a named value implementing fmt.Stringer to the context builder.
func (b *Context) Stg(name string, value fmt.Stringer) *Context {
	b.fields = append(b.fields, contextTuple{
		name: name,
		kind: tupleKindString,
		str:  value.String(),
	})
	return b
}

// Strs adds a named slice of strings value to the context builder.
func (b *Context) Strs(name string, value []string) *Context {
	b.fields = append(b.fields, contextTuple{
		name:  name,
		kind:  tupleKindStrings,
		value: value,
	})
	return b
}

// Bytes adds a named slice of bytes value to the context builder.
//
// Attention:
//
// This operation can be computationally heavy,
// as it evaluates whether the sequence can be represented
// as a string, after which it is saved in either string or object format.
func (b *Context) Bytes(name string, value []byte) *Context {
	tuple := contextTuple{
		name: name,
	}
	if isPrintableStringWithSpaces(value) {
		tuple.kind = tupleKindString
		tuple.str = string(value)
	} else {
		tuple.kind = tupleKindAny
		tuple.value = value
	}

	b.fields = append(b.fields, tuple)
	return b
}

// Type adds a type name to the context builder.
func (b *Context) Type(name string, typ any) *Context {
	b.fields = append(b.fields, contextTuple{
		name: name,
		kind: tupleKindString,
		str:  fmt.Sprintf("%T", typ),
	})
	return b
}

// Any adds an arbitrary named value to the context builder.
func (b *Context) Any(name string, value any) *Context {
	b.fields = append(b.fields, contextTuple{
		name:  name,
		kind:  tupleKindAny,
		value: value,
	})
	return b
}
