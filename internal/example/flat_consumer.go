package main

import (
	"go/token"
	"log/slog"

	"github.com/sirkon/errors"
)

// Implementation of [errors.ErrorContextConsumer] for logging with [slog.Logger].
type flatConsumer struct {
	locs   []annotatedPosition
	fields []any
}

func (c *flatConsumer) NextLink() {}

func (c *flatConsumer) Bool(name string, value bool) {
	c.fields = append(c.fields, slog.Bool(name, value))
}

func (c *flatConsumer) Int(name string, value int) {
	c.fields = append(c.fields, slog.Int(name, value))
}

func (c *flatConsumer) Int8(name string, value int8) {
	c.fields = append(c.fields, slog.Int(name, int(value)))
}

func (c *flatConsumer) Int16(name string, value int16) {
	c.fields = append(c.fields, slog.Int(name, int(value)))
}

func (c *flatConsumer) Int32(name string, value int32) {
	c.fields = append(c.fields, slog.Int(name, int(value)))
}

func (c *flatConsumer) Int64(name string, value int64) {
	c.fields = append(c.fields, slog.Int64(name, value))
}

func (c *flatConsumer) Uint(name string, value uint) {
	c.fields = append(c.fields, slog.Uint64(name, uint64(value)))
}

func (c *flatConsumer) Uint8(name string, value uint8) {
	c.fields = append(c.fields, slog.Uint64(name, uint64(value)))
}

func (c *flatConsumer) Uint16(name string, value uint16) {
	c.fields = append(c.fields, slog.Uint64(name, uint64(value)))
}

func (c *flatConsumer) Uint32(name string, value uint32) {
	c.fields = append(c.fields, slog.Uint64(name, uint64(value)))
}

func (c *flatConsumer) Uint64(name string, value uint64) {
	c.fields = append(c.fields, slog.Uint64(name, value))
}

func (c *flatConsumer) Flt32(name string, value float32) {
	c.fields = append(c.fields, slog.Float64(name, float64(value)))
}

func (c *flatConsumer) Flt64(name string, value float64) {
	c.fields = append(c.fields, slog.Float64(name, value))
}

func (c *flatConsumer) Str(name string, value string) {
	c.fields = append(c.fields, slog.String(name, value))
}

func (c *flatConsumer) Any(name string, value any) {
	c.fields = append(c.fields, slog.Any(name, value))
}

func (c *flatConsumer) SetLinkInfo(loc token.Position, descr errors.ErrorChainLinkDescriptor) {
	var text string
	switch v := descr.(type) {
	case errors.ErrorChainLinkNew:
		text = "NEW: " + string(v)
	case errors.ErrorChainLinkWrap:
		text = "WRAP: " + string(v)
	case errors.ErrorChainLinkContext:
		text = "CTX"
	}
	c.locs = append(c.locs, annotatedPosition{loc, text})
}

var (
	_ errors.ErrorContextConsumer = new(flatConsumer)
)
