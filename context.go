package errors

import (
	"fmt"
	"runtime"
)

// Context возврат контекста ошибки
func Context() ContextBuilder {
	return &errContext{}
}

var _ error = &errContext{}
var _ ContextBuilder = &errContext{}

type location struct {
	file string
	line int
}

type errContext struct {
	items []contextItem
	loc   *location
}

func (c *errContext) Error() string { return "" }

type contextItem struct {
	name  string
	value itemValue
}

type itemValue interface {
	reportItem(name string, builder ContextReportBuilder)
}

// ContextBuilder заполнение контекста информацией и функции завершения построения
type ContextBuilder interface {
	New(msg string) error
	Newf(format string, a ...interface{}) error
	Wrap(err error, msg string) error
	Wrapf(err error, format string, a ...interface{}) error

	Int(name string, value int) ContextBuilder
	Int8(name string, value int8) ContextBuilder
	Int16(name string, value int16) ContextBuilder
	Int32(name string, value int32) ContextBuilder
	Int64(name string, value int64) ContextBuilder
	Uint(name string, value uint) ContextBuilder
	Uint8(name string, value uint8) ContextBuilder
	Uint16(name string, value uint16) ContextBuilder
	Uint32(name string, value uint32) ContextBuilder
	Uint64(name string, value uint64) ContextBuilder
	Float32(name string, value float32) ContextBuilder
	Float64(name string, value float64) ContextBuilder
	String(name string, value string) ContextBuilder
	Strings(name string, values []string) ContextBuilder
	Any(name string, value interface{}) ContextBuilder
	Line() ContextBuilder

	Report(dest ContextReportBuilder)
}

// реализуем контекстом ContextBuilder

func (c *errContext) New(msg string) error {
	return &wrappedError{
		err: New(msg),
		ctx: c,
	}
}

func (c *errContext) Newf(format string, a ...interface{}) error {
	return &wrappedError{
		err: Newf(format, a...),
		ctx: c,
	}
}

func (c *errContext) Wrap(err error, msg string) error {
	var cc *errContext
	if As(err, &cc) {
		// контекст уже существует. Добавляем собранные данные
		cc.items = append(cc.items, c.items...)
		if c.loc != nil {
			// заменяем новой локацией
			cc.loc = c.loc
		}
		return Wrap(err, msg)
	}
	res := Wrap(err, msg).(*wrappedError)
	res.ctx = c
	return res
}

func (c *errContext) Wrapf(err error, format string, a ...interface{}) error {
	return c.Wrap(err, fmt.Sprintf(format, a...))
}

func (c *errContext) Int(name string, value int) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: intValue(value),
	})
	return c
}

func (c *errContext) Uint(name string, value uint) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: uintValue(value),
	})
	return c
}

func (c *errContext) Int8(name string, value int8) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: int8Value(value),
	})
	return c
}

func (c *errContext) Int16(name string, value int16) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: int16Value(value),
	})
	return c
}

func (c *errContext) Int32(name string, value int32) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: int32Value(value),
	})
	return c
}

func (c *errContext) Int64(name string, value int64) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: int64Value(value),
	})
	return c
}

func (c *errContext) Uint8(name string, value uint8) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: uint8Value(value),
	})
	return c
}

func (c *errContext) Uint16(name string, value uint16) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: uint16Value(value),
	})
	return c
}

func (c *errContext) Uint32(name string, value uint32) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: uint32Value(value),
	})
	return c
}

func (c *errContext) Uint64(name string, value uint64) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: uint64Value(value),
	})
	return c
}

func (c *errContext) Float32(name string, value float32) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: float32Value(value),
	})
	return c
}

func (c *errContext) Float64(name string, value float64) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: float64Value(value),
	})
	return c
}

func (c *errContext) String(name string, value string) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: stringValue(value),
	})
	return c
}

func (c *errContext) Strings(name string, values []string) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: stringsValue(values),
	})
	return c
}

func (c *errContext) Any(name string, value interface{}) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: anyValue{val: value},
	})
	return c
}

func (c *errContext) Line() ContextBuilder {
	_, fn, line, _ := runtime.Caller(1)
	if c.loc == nil {
		c.loc = &location{}
	}
	c.loc.file = fn
	c.loc.line = line
	return c
}

func (c *errContext) Report(dest ContextReportBuilder) {
	for _, item := range c.items {
		item.value.reportItem(item.name, dest)
	}
}
