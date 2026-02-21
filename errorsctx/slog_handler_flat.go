package errorsctx

import (
	"context"
	"log/slog"

	errors "awesome-errors"
)

// SLogHandlerFlat handler for a flat view of an error context.
type SLogHandlerFlat struct {
	handler slog.Handler
}

func NewSLogHandlerFlat(handler slog.Handler) *SLogHandlerFlat {
	return &SLogHandlerFlat{handler: handler}
}

func (h *SLogHandlerFlat) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *SLogHandlerFlat) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.handler.WithAttrs(attrs)
}

func (h *SLogHandlerFlat) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}

func (h *SLogHandlerFlat) Handle(ctx context.Context, r slog.Record) error {
	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

	r.Attrs(func(a slog.Attr) bool {
		e, ok := a.Value.Any().(error)
		if !ok {
			newRecord.AddAttrs(a)
			return true
		}

		err, ok := e.(*errors.Error)
		if !ok {
			err, ok = errors.AsType[*errors.Error](e)
			if !ok {
				newRecord.AddAttrs(a)
				return true
			}
		}

		if a.Key == "" || a.Key == "!BADKEY" {
			a.Key = "err"
		}

		// Add error message under a key.
		newRecord.AddAttrs(slog.String(a.Key, err.Error()))

		// Add context tree as @key.
		ctx := errors.SLogFlatContext(err)
		newRecord.AddAttrs(slog.GroupAttrs("@"+a.Key, ctx...))
		return true

	})

	return h.handler.Handle(ctx, newRecord)
}
