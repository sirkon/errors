package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"math"
	"os"

	errors "awesome-errors"
	"awesome-errors/errorsctx"
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
