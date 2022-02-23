package errors

import "fmt"

func (e Error) Bool(name string, value bool) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Int(name string, value int) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}
func (e Error) Int8(name string, value int8) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Int16(name string, value int16) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Int32(name string, value int32) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Int64(name string, value int64) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Uint(name string, value uint) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Uint8(name string, value uint8) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Uint16(name string, value uint16) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Uint32(name string, value uint32) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Uint64(name string, value uint64) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Float32(name string, value float32) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Float64(name string, value float64) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Str(name string, value string) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Stg(name string, value fmt.Stringer) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Strs(name string, value []string) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}

func (e Error) Any(name string, value interface{}) Error {
	return Error{
		msg: e.msg,
		err: e.err,
		ctx: append(e.ctx, contextTuple{
			name:  name,
			value: value,
		}),
	}
}
