package errors

import (
	"fmt"
	"log/slog"
	"unicode"
	"unicode/utf8"
)

func (e *Error) Bool(key string, value bool) *Error {
	e.pushValue(key, slog.BoolValue(value), errorAttrKindBool)
	return e
}

func (e *Error) Int(key string, value int) *Error {
	e.pushValue(key, slog.Int64Value(int64(value)), errorAttrKindI64)
	return e
}

func (e *Error) I8(key string, value int8) *Error {
	e.pushValue(key, slog.Int64Value(int64(value)), errorAttrKindI64)
	return e
}

func (e *Error) I16(key string, value int16) *Error {
	e.pushValue(key, slog.Int64Value(int64(value)), errorAttrKindI64)
	return e
}

func (e *Error) I32(key string, value int32) *Error {
	e.pushValue(key, slog.Int64Value(int64(value)), errorAttrKindI64)
	return e
}

func (e *Error) I64(key string, value int64) *Error {
	e.pushValue(key, slog.Int64Value(value), errorAttrKindI64)
	return e
}

func (e *Error) Uint(key string, value uint) *Error {
	e.pushValue(key, slog.Uint64Value(uint64(value)), errorAttrKindU64)
	return e
}

func (e *Error) U8(key string, value uint8) *Error {
	e.pushValue(key, slog.Uint64Value(uint64(value)), errorAttrKindU64)
	return e
}

func (e *Error) U16(key string, value uint16) *Error {
	e.pushValue(key, slog.Uint64Value(uint64(value)), errorAttrKindU64)
	return e
}

func (e *Error) U32(key string, value uint32) *Error {
	e.pushValue(key, slog.Uint64Value(uint64(value)), errorAttrKindU64)
	return e
}

func (e *Error) U64(key string, value uint64) *Error {
	e.pushValue(key, slog.Uint64Value(value), errorAttrKindU64)
	return e
}

func (e *Error) F32(key string, value float32) *Error {
	e.pushValue(key, slog.Float64Value(float64(value)), errorAttrKindF64)
	return e
}

func (e *Error) F64(key string, value float64) *Error {
	e.pushValue(key, slog.Float64Value(value), errorAttrKindF64)
	return e
}

func (e *Error) Str(key, value string) *Error {
	e.pushValue(key, slog.StringValue(value), errorAttrKindStr)
	return e
}

func (e *Error) Stg(key string, value fmt.Stringer) *Error {
	e.pushValue(key, slog.StringValue(value.String()), errorAttrKindStr)
	return e
}

func (e *Error) Strs(key string, value []string) *Error {
	e.pushValue(key, slog.AnyValue(value), errorAttrKindAny)
	return e
}

func (e *Error) Bytes(key string, value []byte) *Error {
	if isPrintableStringWithSpaces(value) {
		e.pushValue(key, slog.StringValue(string(value)), errorAttrKindStr)
		return e
	}

	e.pushValue(key, slog.AnyValue(value), errorAttrKindAny)
	return e
}

func (e *Error) Any(key string, value any) *Error {
	e.pushValue(key, slog.AnyValue(value), errorAttrKindAny)
	return e
}

func (e *Error) pushValue(key string, value slog.Value, kind errorAttrKind) {
	e.attrs = append(e.attrs, errorAttr{
		kind:  kind,
		key:   key,
		value: value,
	})
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
