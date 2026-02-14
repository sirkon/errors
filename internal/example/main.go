package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"

	"github.com/sirkon/errors"
)

func main() {
	errors.InsertLocations()

	// Create a new error.
	err := errors.New("failed to do something").
		Int("int-value", 13).
		Str("string-value", "world")

	// Annotate the created error.
	err = errors.Wrap(err, "ask to do something").Bool("insert-locations", true)

	// Add additional structured context.
	err = errors.Just(err).Flt64("pi", math.Pi).Flt64("e", math.Exp(1))

	// Output "prettily". In general, it would be correct to do this at the
	// slog.Handler level, but in our case this is sufficient.
	//
	// At the same time, prepare to save the output to the OUTPUT.md file to avoid manually copying
	// the new output after any changes.
	w := &prettyJSONWriter{
		collect: bytes.NewBuffer(nil),
	}

	// Reconfigure slog for convenient human-readable output.
	slog.SetDefault(slog.New(errors.NewSLogErrorContextGrouppedHandler(slog.NewJSONHandler(
		w,
		&slog.HandlerOptions{
			AddSource:   true,
			Level:       nil,
			ReplaceAttr: nil,
		},
	))))
	slog.Info("hello world", err, "name", "value")
	md1 := w.collect.String()

	// Generate OUTPUT.md content
	w.collect.Reset()
	slog.SetDefault(slog.New(errors.NewSLogErrorContextFlatHandler(slog.NewJSONHandler(
		w,
		&slog.HandlerOptions{
			AddSource:   true,
			Level:       nil,
			ReplaceAttr: nil,
		},
	))))
	slog.Info("hello world", err, "name", "value")
	md2 := w.collect.String()

	mdFormat := `# Logging output example.

## Error context grouped by places where it was added.
%s

## Flat structure, error context is shown all together.
%s
`
	md := fmt.Sprintf(mdFormat, mdJSON(md1), mdJSON(md2))

	const outputExample = "OUTPUT.md"
	if err := os.WriteFile("./internal/example/"+outputExample, []byte(md), 0644); err != nil {
		fmt.Printf("save output example into %s.md: %s\n", outputExample, err)
		os.Exit(1)
	}
}

func mdJSON(v string) string {
	return fmt.Sprintf("```json\n%s```", v)
}

// Special implementation of [io.Writer] for outputting JSON in a human-readable format.
type prettyJSONWriter struct {
	collect *bytes.Buffer
}

func (w *prettyJSONWriter) Write(p []byte) (n int, err error) {
	var dst bytes.Buffer
	if err := json.Indent(&dst, p, "", "  "); err != nil {
		return 0, errors.Wrap(err, "format JSON log line")
	}

	mw := io.MultiWriter(os.Stdout, w.collect)
	written, err := io.Copy(mw, &dst)
	if err != nil {
		return 0, errors.Wrap(err, "write log line into STDOUT")
	}

	return int(written), nil
}
