package errors

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

// Logger возвращаем логгер аннотированный, возможно, значениями из контекста ошибки (если он заполнен)
func Logger(logger zerolog.Logger, err error) *zerolog.Logger {
	b := &zerologReportBuilder{
		ctx: logger.With(),
	}
	b.Err(err)
	var c *errContext
	if As(err, &c) {
		if c.loc != nil {
			b.String("err-pos", fmt.Sprintf("%s:%d", c.loc.file, c.loc.line))
		}
		c.Report(b)
	}
	res := b.Logger(logger)
	return &res
}

// CtxLogger вытаскивает логгер из контекста и запускает с ним Logger
func CtxLogger(ctx context.Context, err error) *zerolog.Logger {
	return Logger(*zerolog.Ctx(ctx), err)
}
