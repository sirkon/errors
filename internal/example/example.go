package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"runtime/debug"

	"github.com/sirkon/errors"
	"github.com/sirkon/errors/errorsctx"
)

func main() {
	errors.InsertLocations()
	err := errors.New("this is an error").
		Bytes("bytes", []byte{1, 2, 3}).
		Bytes("text-bytes", []byte("Hello World!"))
	err = errors.Wrap(err, "check error").
		Int("count", 333).
		Bool("is-wrap-layer", true)
	err = errors.Just(err).
		F64("pi", math.Pi).
		F64("e", math.E)

	logger := slog.New(
		errorsctx.NewSlogPrettyRenderer(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
			true,
			-1,
		),
	)

	logger.Error("pure foreign error", io.EOF)
	logger.Error("log error with just layers", (errors.Wrap(io.EOF, "wrap")))
	logger.Error("log error with tree structured context", err)
	logger.Error("log marked error with tree", errors.Spec(err, new(0)))
	logger.Error("log marked foreign error with tree beneath",
		errors.Spec(fmt.Errorf("foreign wrap: %w", err), new(0)),
	)
	logger.Info("simple info")
	logger.Info("simple info with ctx2", slog.Int("count", 42), slog.String("key", "value"))
	logger.Info("simple info with ctx4",
		slog.Int("count", 42),
		slog.String("key", "value"),
		slog.String("key", "value"),
		slog.String("key", "value"),
	)
	logger.Info("simple stack", slog.String("stack", string(debug.Stack())))
	logger.Info(
		"with internal json",
		slog.Any("obj", map[string]any{
			"foo": "bar",
			"data": map[string]int{
				"k":  1,
				"k2": 2,
			},
		}),
		errorsctx.ForceTree(),
	)
	logger.Info("with internal text json tree", slog.String("obj", `{"foo": "bar"}`))
	logger.Info("with internal text json array", slog.String("obj", `[1,2,3]`))

	// logger = slog.New(errorsctx.NewSLogHandlerFlat(
	// 	slog.NewJSONHandler(&fancyJSONWriter{}, &slog.HandlerOptions{}),
	// ))
	// logger.Error("log error with flat structured context", err)
}

type fancyJSONWriter struct{}

func (*fancyJSONWriter) Write(p []byte) (n int, err error) {
	var buf bytes.Buffer
	_ = json.Indent(&buf, p, "", "  ")
	_, _ = os.Stderr.Write(buf.Bytes())
	return len(p), nil
}
