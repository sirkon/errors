package main

import (
	"go/token"
	"log/slog"

	"github.com/sirkon/errors"
)

type basicConsumer struct {
	loc    token.Position
	descr  errors.ErrorChainLinkDescriptor
	text   string
	fields []any
}

// Implementation of [errors.ErrorContextConsumer] for logging with [slog.Logger] with separation
// of context by the place where it was added.
type groupedConsumer struct {
	consumers []*basicConsumer
	consumer  *basicConsumer
}

func newGroupedConsumer() *groupedConsumer {
	return &groupedConsumer{}
}

func (c *groupedConsumer) push(field any) {
	cons := c.consumers[len(c.consumers)-1]
	cons.fields = append(cons.fields, field)
}

func (c *groupedConsumer) NextLink() {
	cons := &basicConsumer{}
	c.consumers = append(c.consumers, cons)
	c.consumer = cons
}

func (c *groupedConsumer) Bool(name string, value bool) {
	c.push(slog.Bool(name, value))
}

func (c *groupedConsumer) Int(name string, value int) {
	c.push(slog.Int(name, value))
}

func (c *groupedConsumer) Int8(name string, value int8) {
	c.push(slog.Int(name, int(value)))
}

func (c *groupedConsumer) Int16(name string, value int16) {
	c.push(slog.Int(name, int(value)))
}

func (c *groupedConsumer) Int32(name string, value int32) {
	c.push(slog.Int(name, int(value)))
}

func (c *groupedConsumer) Int64(name string, value int64) {
	c.push(slog.Int64(name, value))
}

func (c *groupedConsumer) Uint(name string, value uint) {
	c.push(slog.Uint64(name, uint64(value)))
}

func (c *groupedConsumer) Uint8(name string, value uint8) {
	c.push(slog.Uint64(name, uint64(value)))
}

func (c *groupedConsumer) Uint16(name string, value uint16) {
	c.push(slog.Uint64(name, uint64(value)))
}

func (c *groupedConsumer) Uint32(name string, value uint32) {
	c.push(slog.Uint64(name, uint64(value)))
}

func (c *groupedConsumer) Uint64(name string, value uint64) {
	c.push(slog.Uint64(name, value))
}

func (c *groupedConsumer) Flt32(name string, value float32) {
	c.push(slog.Float64(name, float64(value)))
}

func (c *groupedConsumer) Flt64(name string, value float64) {
	c.push(slog.Float64(name, value))
}

func (c *groupedConsumer) Str(name string, value string) {
	c.push(slog.String(name, value))
}

func (c *groupedConsumer) Any(name string, value any) {
	c.push(slog.Any(name, value))
}

func (c *groupedConsumer) SetLinkInfo(loc token.Position, descr errors.ErrorChainLinkDescriptor) {
	switch v := descr.(type) {
	case errors.ErrorChainLinkNew:
		c.consumer.text = "NEW: " + string(v)
	case errors.ErrorChainLinkWrap:
		c.consumer.text = "WRAP: " + string(v)
	case errors.ErrorChainLinkContext:
		c.consumer.text = "CTX"
	}
	c.consumer.descr = descr
	c.consumer.loc = loc
}

var (
	_ errors.ErrorContextConsumer = new(groupedConsumer)
)
