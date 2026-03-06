package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"os"

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

	logger := slog.New(errorsctx.NewSLogHandlerTree(
		slog.NewJSONHandler(&fancyJSONWriter{}, &slog.HandlerOptions{}),
	))

	logger.Error("log error with tree structured context", err)
	logger.Error("log marked error with tree", errors.Mark(err, new(0)))
	logger.Error("log marked foreign error with tree beneath",
		errors.Mark(fmt.Errorf("foreign wrap: %w", err), new(0)),
	)

	logger = slog.New(errorsctx.NewSLogHandlerFlat(
		slog.NewJSONHandler(&fancyJSONWriter{}, &slog.HandlerOptions{}),
	))
	logger.Error("log error with flat structured context", err)
}

type fancyJSONWriter struct{}

func (*fancyJSONWriter) Write(p []byte) (n int, err error) {
	var buf bytes.Buffer
	_ = json.Indent(&buf, p, "", "  ")
	_, _ = os.Stderr.Write(buf.Bytes())
	return len(p), nil
}
