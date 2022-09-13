package errors

import (
	"fmt"
	"runtime"
	"strconv"
)

// Pfx adds (replaces) prefix in the rest of the chain.
func (e Error) Pfx(prefix string) Error {
	if prefix == "" {
		return e
	}

	e.ctxPrefix = prefix + "-"
	return e
}

// Bool adds boolean named value into the context
func (e Error) Bool(name string, value bool) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Int adds int named value into the context
func (e Error) Int(name string, value int) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Int8 adds int8 named value into the context
func (e Error) Int8(name string, value int8) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Int16 adds int16 named value into the context
func (e Error) Int16(name string, value int16) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Int32 adds int32 named value into the context
func (e Error) Int32(name string, value int32) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Int64 adds int64 named value into the context
func (e Error) Int64(name string, value int64) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Uint adds uint named value into the context
func (e Error) Uint(name string, value uint) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Uint8 adds uint8 named value into the context
func (e Error) Uint8(name string, value uint8) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Uint16 adds uint16 named value into the context
func (e Error) Uint16(name string, value uint16) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Uint32 adds uint32 named value into the context
func (e Error) Uint32(name string, value uint32) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Uint64 adds uint64 named value into the context
func (e Error) Uint64(name string, value uint64) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Float32 adds float32 named value into the context
func (e Error) Float32(name string, value float32) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Float64 adds float64 named value into the context
func (e Error) Float64(name string, value float64) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Str adds string named value into the context
func (e Error) Str(name string, value string) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Stg adds named value of the given stringer into the context
func (e Error) Stg(name string, value fmt.Stringer) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Strs adds named slice of strings value into the context
func (e Error) Strs(name string, value []string) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

// Any adds some named value into the context
func (e Error) Any(name string, value interface{}) Error {
	return Error{
		msg:       e.msg,
		err:       e.err,
		ctxPrefix: e.ctxPrefix,
		ctx: append(e.ctx, contextTuple{
			name:  e.ctxPrefix + name,
			value: value,
		}),
	}
}

const (
	locationName = "location"
)

// Loc adds error location into the context
func (e Error) Loc(depth int) Error {
	_, fn, line, _ := runtime.Caller(1 + depth)
	return e.Str(locationName, fn+":"+strconv.Itoa(line))
}
