package errorsctx

import (
	"context"
	"log/slog"

	"github.com/sirkon/errors"
)

// SLogHandlerTree handler for a tree view of an error context.
type SLogHandlerTree struct {
	handler slog.Handler
}

// NewSLogHandlerTree creates handler [SLogHandlerTree].
func NewSLogHandlerTree(handler slog.Handler) *SLogHandlerTree {
	return &SLogHandlerTree{
		handler: handler,
	}
}

func (h *SLogHandlerTree) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *SLogHandlerTree) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.handler.WithAttrs(attrs)
}

func (h *SLogHandlerTree) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}

// Handle handles errors.
func (h *SLogHandlerTree) Handle(ctx context.Context, r slog.Record) error {
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
		treeContext := errors.SLogTreeContext(err)
		newRecord.AddAttrs(slog.GroupAttrs("@"+a.Key, treeContext...))
		return true

	})

	return h.handler.Handle(ctx, newRecord)
}
