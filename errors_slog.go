package errors

import (
	"context"
	"go/token"
	"log/slog"
)

func NewSLogErrorContextGrouppedHandler(handler slog.Handler) *SLogErrorContextGrouppedHandler {
	return &SLogErrorContextGrouppedHandler{
		Handler: handler,
	}
}

func NewSLogErrorContextFlatHandler(handler slog.Handler) *SLogERrorContextFlatHandler {
	return &SLogERrorContextFlatHandler{
		Handler: handler,
	}
}

// SLogErrorContextGrouppedHandler to handle these errors with slog.Logger splitting errors by processing stages.
type SLogErrorContextGrouppedHandler struct {
	slog.Handler
}

const contextLengthPrediction = 4

// Handle handles errors.
func (h *SLogErrorContextGrouppedHandler) Handle(ctx context.Context, r slog.Record) error {
	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

	r.Attrs(func(a slog.Attr) bool {
		if err, ok := a.Value.Any().(*Error); ok {
			if a.Key == "" || a.Key == "!BADKEY" {
				a.Key = "err"
			}

			// Add error message as key
			newRecord.AddAttrs(slog.String(a.Key, err.Error()))

			// Add context tree as @key.
			dlvr := GetContextDeliverer(err)
			cons := newGroupedConsumer()
			dlvr.Deliver(cons)
			errCtxFields := make([]any, 0, contextLengthPrediction)

			for _, c := range cons.consumers {
				fs := make([]any, 0, contextLengthPrediction)
				if c.loc.IsValid() {
					fs = append(fs, slog.String("@location", c.loc.String()))
				}
				for _, f := range c.fields {
					fs = append(fs, f)
				}
				errCtxFields = append(errCtxFields, slog.Group(c.text, fs...))
			}
			newRecord.AddAttrs(slog.Group("@"+a.Key, errCtxFields...))
			return true
		}
		newRecord.AddAttrs(a)
		return true
	})

	return h.Handler.Handle(ctx, newRecord)
}

func (h *SLogErrorContextGrouppedHandler) WithGroup(name string) slog.Handler {
	return &SLogErrorContextGrouppedHandler{Handler: h.Handler.WithGroup(name)}
}

func (h *SLogErrorContextGrouppedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SLogErrorContextGrouppedHandler{Handler: h.Handler.WithAttrs(attrs)}
}

// SLogERrorContextFlatHandler to handle these errors with slog.Logger splitting errors by processing stages.
type SLogERrorContextFlatHandler struct {
	slog.Handler
}

// Handle handles errors.
func (h *SLogERrorContextFlatHandler) Handle(ctx context.Context, r slog.Record) error {
	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

	r.Attrs(func(a slog.Attr) bool {
		if err, ok := a.Value.Any().(*Error); ok {
			if a.Key == "" || a.Key == "!BADKEY" {
				a.Key = "err"
			}

			// Add error message as key
			newRecord.AddAttrs(slog.String(a.Key, err.Error()))

			// Add context tree as @key.
			dlvr := GetContextDeliverer(err)
			cons := &slogFlatConsumer{}
			dlvr.Deliver(cons)
			l := len(cons.fields)
			if len(cons.locs) > 0 {
				l++
			}
			fields := make([]any, l)
			copy(fields, cons.fields)
			if len(cons.locs) > 0 {
				fields[len(fields)-1] = slog.Group("@locations", cons.locs...)
			}
			newRecord.AddAttrs(slog.Group("@"+a.Key, fields...))
			return true
		}
		newRecord.AddAttrs(a)
		return true
	})

	return h.Handler.Handle(ctx, newRecord)
}

func (h *SLogERrorContextFlatHandler) WithGroup(name string) slog.Handler {
	return &SLogERrorContextFlatHandler{Handler: h.Handler.WithGroup(name)}
}

func (h *SLogERrorContextFlatHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SLogERrorContextFlatHandler{Handler: h.Handler.WithAttrs(attrs)}
}

type basicConsumer struct {
	loc    token.Position
	descr  ErrorChainLinkDescriptor
	text   string
	fields []slog.Attr
}

// [errors.ErrorContextConsumer] implementation to log with [slog.Logger] with
// processing level separation.
type slogGrouppedConsumer struct {
	consumers []*basicConsumer
	consumer  *basicConsumer
}

func newGroupedConsumer() *slogGrouppedConsumer {
	return &slogGrouppedConsumer{}
}

func (c *slogGrouppedConsumer) push(field slog.Attr) {
	cons := c.consumers[len(c.consumers)-1]
	cons.fields = append(cons.fields, field)
}

func (c *slogGrouppedConsumer) NextLink() {
	cons := &basicConsumer{}
	c.consumers = append(c.consumers, cons)
	c.consumer = cons
}

func (c *slogGrouppedConsumer) Bool(name string, value bool) {
	c.push(slog.Bool(name, value))
}

func (c *slogGrouppedConsumer) Int(name string, value int) {
	c.push(slog.Int(name, value))
}

func (c *slogGrouppedConsumer) Int8(name string, value int8) {
	c.push(slog.Int(name, int(value)))
}

func (c *slogGrouppedConsumer) Int16(name string, value int16) {
	c.push(slog.Int(name, int(value)))
}

func (c *slogGrouppedConsumer) Int32(name string, value int32) {
	c.push(slog.Int(name, int(value)))
}

func (c *slogGrouppedConsumer) Int64(name string, value int64) {
	c.push(slog.Int64(name, value))
}

func (c *slogGrouppedConsumer) Uint(name string, value uint) {
	c.push(slog.Uint64(name, uint64(value)))
}

func (c *slogGrouppedConsumer) Uint8(name string, value uint8) {
	c.push(slog.Uint64(name, uint64(value)))
}

func (c *slogGrouppedConsumer) Uint16(name string, value uint16) {
	c.push(slog.Uint64(name, uint64(value)))
}

func (c *slogGrouppedConsumer) Uint32(name string, value uint32) {
	c.push(slog.Uint64(name, uint64(value)))
}

func (c *slogGrouppedConsumer) Uint64(name string, value uint64) {
	c.push(slog.Uint64(name, value))
}

func (c *slogGrouppedConsumer) Flt32(name string, value float32) {
	c.push(slog.Float64(name, float64(value)))
}

func (c *slogGrouppedConsumer) Flt64(name string, value float64) {
	c.push(slog.Float64(name, value))
}

func (c *slogGrouppedConsumer) Str(name string, value string) {
	c.push(slog.String(name, value))
}

func (c *slogGrouppedConsumer) Any(name string, value any) {
	c.push(slog.Any(name, value))
}

func (c *slogGrouppedConsumer) SetLinkInfo(loc token.Position, descr ErrorChainLinkDescriptor) {
	switch v := descr.(type) {
	case ErrorChainLinkNew:
		c.consumer.text = "NEW: " + string(v)
	case ErrorChainLinkWrap:
		c.consumer.text = "WRAP: " + string(v)
	case ErrorChainLinkContext:
		c.consumer.text = "CTX"
	}
	c.consumer.descr = descr
	c.consumer.loc = loc
}

// Implementation of [errors.ErrorContextConsumer] for logging with [slog.Logger].
type slogFlatConsumer struct {
	locs   []any
	fields []any
}

func (c *slogFlatConsumer) NextLink() {}

func (c *slogFlatConsumer) Bool(name string, value bool) {
	c.fields = append(c.fields, slog.Bool(name, value))
}

func (c *slogFlatConsumer) Int(name string, value int) {
	c.fields = append(c.fields, slog.Int(name, value))
}

func (c *slogFlatConsumer) Int8(name string, value int8) {
	c.fields = append(c.fields, slog.Int(name, int(value)))
}

func (c *slogFlatConsumer) Int16(name string, value int16) {
	c.fields = append(c.fields, slog.Int(name, int(value)))
}

func (c *slogFlatConsumer) Int32(name string, value int32) {
	c.fields = append(c.fields, slog.Int(name, int(value)))
}

func (c *slogFlatConsumer) Int64(name string, value int64) {
	c.fields = append(c.fields, slog.Int64(name, value))
}

func (c *slogFlatConsumer) Uint(name string, value uint) {
	c.fields = append(c.fields, slog.Uint64(name, uint64(value)))
}

func (c *slogFlatConsumer) Uint8(name string, value uint8) {
	c.fields = append(c.fields, slog.Uint64(name, uint64(value)))
}

func (c *slogFlatConsumer) Uint16(name string, value uint16) {
	c.fields = append(c.fields, slog.Uint64(name, uint64(value)))
}

func (c *slogFlatConsumer) Uint32(name string, value uint32) {
	c.fields = append(c.fields, slog.Uint64(name, uint64(value)))
}

func (c *slogFlatConsumer) Uint64(name string, value uint64) {
	c.fields = append(c.fields, slog.Uint64(name, value))
}

func (c *slogFlatConsumer) Flt32(name string, value float32) {
	c.fields = append(c.fields, slog.Float64(name, float64(value)))
}

func (c *slogFlatConsumer) Flt64(name string, value float64) {
	c.fields = append(c.fields, slog.Float64(name, value))
}

func (c *slogFlatConsumer) Str(name string, value string) {
	c.fields = append(c.fields, slog.String(name, value))
}

func (c *slogFlatConsumer) Any(name string, value any) {
	c.fields = append(c.fields, slog.Any(name, value))
}

func (c *slogFlatConsumer) SetLinkInfo(loc token.Position, descr ErrorChainLinkDescriptor) {
	var text string
	switch v := descr.(type) {
	case ErrorChainLinkNew:
		text = "NEW: " + string(v)
	case ErrorChainLinkWrap:
		text = "WRAP: " + string(v)
	case ErrorChainLinkContext:
		text = "CTX"
	}
	if !loc.IsValid() {
		return
	}
	c.locs = append(c.locs, slog.String(loc.String(), text))
}

var (
	_ ErrorContextConsumer = new(slogFlatConsumer)
	_ ErrorContextConsumer = new(slogGrouppedConsumer)
)
