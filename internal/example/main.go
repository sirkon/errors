package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/token"
	"io"
	"log/slog"
	"math"
	"os"

	"github.com/sirkon/errors"
)

func main() {
	// Output "prettily". In general, it would be correct to do this at the
	// slog.Handler level, but in our case this is sufficient.
	//
	// At the same time, prepare to save the output to the OUTPUT.md file to avoid manually copying
	// the new output after any changes.
	w := &prettyJSONWriter{
		collect: bytes.NewBuffer(nil),
	}

	// Reconfigure slog for convenient human-readable output.
	slog.SetDefault(slog.New(slog.NewJSONHandler(
		w,
		&slog.HandlerOptions{
			AddSource:   true,
			Level:       nil,
			ReplaceAttr: nil,
		},
	)))
	errors.InsertLocations()

	//
	// Create an error for logging.

	// Create a new error.
	err := errors.New("failed to do something").
		Int("int-value", 13).
		Str("string-value", "world")

	// Annotate the created error.
	err = errors.Wrap(err, "ask to do something").Bool("insert-locations", true)

	// Add additional structured context.
	err = errors.Just(err).Flt64("pi", math.Pi).Flt64("e", math.Exp(1))

	// Log with logging context and error context.
	LogFlat(
		"logging test with flat output",
		slog.Int("int-value", 12),
		slog.String("string-value", "hello"),
		errors.Log(err), // Same as slog.Any("err", err).
		slog.Float64("x", 1.5),
		slog.Any("err-naked", errors.Const("naked error message")),
	)

	// Save the first output example for OUTPUT.md
	md1 := w.collect.String()
	w.collect.Reset()

	// Second example, error context will be grouped by processing location.
	LogGrouped(
		"logging test with error context grouped around places it was added",
		slog.Bool("grouped-structure", true),
		errors.Log(err),
		slog.String("greeting", "Hello, World!"),
	)

	// Generate OUTPUT.md content
	md2 := w.collect.String()
	mdFormat := `# Logging output example.

## Flat structure, error context is shown all together.
%s

## Error context grouped by places where it was added.
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

// LogFlat function for logging to slog at Info level with extraction of context from error objects
// received in fields. Uses "flat" error context structure.
func LogFlat(msg string, fields ...slog.Attr) {
	attrs := feedFieldsWithFlatErrorContext(fields)
	slog.Info(msg, attrs...)
}

// LogGrouped similar to LogFlat, but uses error structure separated by processing location.
func LogGrouped(msg string, fields ...slog.Attr) {
	attrs := feedFieldsWithGroupedErrorContext(fields)
	slog.Info(msg, attrs...)
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

const (
	locationsLoggingKey = "@locations"
)

func errContextKey(errname string) string {
	return "@" + errname
}

type Logger struct {
}

// This function modifies the list of fields, adding error context if present.
func feedFieldsWithFlatErrorContext(fields []slog.Attr) []any {
	const contextLengthPrediction = 4

	attrs := make([]any, 0, len(fields)+contextLengthPrediction)
	for _, field := range fields {
		// Save all original fields preserving their order.
		attrs = append(attrs, field)

		// Separately process fields containing errors.
		// Errors, as "complex" types, are contained in Any().
		fv := field.Value.Any()
		if fv == nil {
			continue
		}
		e, ok := fv.(error)
		if !ok {
			continue
		}

		// Try to extract context.
		d := errors.GetContextDeliverer(e)
		if d == nil {
			// This is not our error.
			continue
		}

		// Create a container for the extracted context and fill it with values.
		// Then add the accumulated values to the list of fields for logging
		// as part of the @<err-key> group.
		cons := flatConsumer{
			fields: make([]any, 0, contextLengthPrediction),
		}
		d.Deliver(&cons)

		// Create the content of the @location group to display error processing locations
		// with descriptions of the processing method.
		var locs []any
		for _, loc := range cons.locs {
			if !loc.IsValid() {
				continue
			}
			locs = append(locs, slog.String(loc.String(), loc.Descr))
		}
		if len(locs) > 0 {
			cons.fields = append(cons.fields, slog.Group(locationsLoggingKey, locs...))
		}

		attrs = append(attrs, slog.Group(errContextKey(field.Key), cons.fields...))
	}

	return attrs
}

// This function modifies the list of fields, adding error context if present, grouping it by
// processing location.
func feedFieldsWithGroupedErrorContext(fields []slog.Attr) []any {
	const contextLengthPrediction = 4

	attrs := make([]any, 0, len(fields)+contextLengthPrediction)
	for _, field := range fields {
		// Save all original fields preserving their order.
		attrs = append(attrs, field)

		// Separately process fields containing errors.
		// Errors, as "complex" types, are contained in Any().
		fv := field.Value.Any()
		if fv == nil {
			continue
		}
		e, ok := fv.(error)
		if !ok {
			continue
		}

		// Try to extract context.
		d := errors.GetContextDeliverer(e)
		if d == nil {
			// This is not our error.
			continue
		}

		cons := newGroupedConsumer()
		d.Deliver(cons)

		errCtxFields := make([]any, 0, contextLengthPrediction)
		for _, c := range cons.consumers {
			fs := make([]any, 0, contextLengthPrediction)
			if c.loc.IsValid() {
				fs = append(fs, slog.String("@location", c.loc.String()))
			}
			fs = append(fs, c.fields...)
			errCtxFields = append(errCtxFields, slog.Group(c.text, fs...))
		}

		attrs = append(attrs, slog.Group(errContextKey(field.Key), errCtxFields...))
	}

	return attrs
}

type annotatedPosition struct {
	token.Position
	Descr string
}
