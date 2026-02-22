package errorsctx_test

import (
	"fmt"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirkon/errors"
	"github.com/sirkon/errors/errorsctx"
)

var (
	treeLogger   *slog.Logger
	flatLogger   *slog.Logger
	stdLogger    *slog.Logger
	txtCtxLogger *slog.Logger

	treeFile   *os.File
	flatFile   *os.File
	stdFile    *os.File
	txtCtxFile *os.File
)

func BenchmarkErrorsTree(b *testing.B) {
	b.Cleanup(func() {
		if err := treeFile.Close(); err != nil {
			b.Error("failed to close tree log file:", err)
		}
	})

	for b.Loop() {
		err := errors.New("this is an error").
			Any("bytes", []byte{1, 2, 3}).
			Str("text-bytes", "Hello World!")
		err = errors.Wrap(err, "check error").
			Int("count", 333).
			Bool("is-wrap-layer", true)
		err = errors.Just(err).
			F64("pi", math.Pi).
			F64("e", math.E)

		treeLogger.Error("failed to do something", slog.Any("err", err))
	}
}

func BenchmarkErrorsFlat(b *testing.B) {
	b.Cleanup(func() {
		if err := flatFile.Close(); err != nil {
			b.Error("failed to close flat log file:", err)
		}
	})

	for b.Loop() {
		err := errors.New("this is an error").
			Any("bytes", []byte{1, 2, 3}).
			Str("text-bytes", "Hello World!")
		err = errors.Wrap(err, "check error").
			Int("count", 333).
			Bool("is-wrap-layer", true)
		err = errors.Just(err).
			F64("pi", math.Pi).
			F64("e", math.E)

		flatLogger.Error("failed to do something", slog.Any("err", err))
	}
}

func BenchmarkErrorsStd(b *testing.B) {
	b.Cleanup(func() {
		if err := stdFile.Close(); err != nil {
			b.Error("failed to close std log file:", err)
		}
	})

	for b.Loop() {
		err := fmt.Errorf("this is an error")
		stdLogger.Error(
			"failed to do something 1",
			slog.Any("err", err),
			slog.Any("bytes", []byte{1, 2, 3}),
			slog.String("text-bytes", "Hello World!"),
		)

		stdLogger.Error(
			"failed to check error",
			slog.Any("err", err),
			slog.Int("count", 333),
			slog.Bool("is-wrap-layer", true),
		)
		err = fmt.Errorf("check error: %w", err)

		stdLogger.Error(
			"got an error",
			slog.Any("err", err),
			slog.Float64("pi", math.Pi),
			slog.Float64("e", math.E),
		)
	}
}

func BenchmarkErrorsTxtContext(b *testing.B) {
	b.Cleanup(func() {
		if err := txtCtxFile.Close(); err != nil {
			b.Error("failed to close std log file:", err)
		}
	})

	for b.Loop() {
		err := fmt.Errorf("this is an error bytes[%v] text-bytes[%s]", []byte{1, 2, 3}, "Hello World!")
		err = fmt.Errorf("check error count[%d] is-wrap-layer[%v]: %w", 333, true, err)
		err = fmt.Errorf("context pi[%g] e[%g]: %w", math.Pi, math.E, err)

		txtCtxLogger.Error("failed to do something 1", slog.Any("err", err))
	}
}

func init() {
	var err error
	treeFile, err = os.Create(filepath.Join(os.TempDir(), "errors_tree.log"))
	if err != nil {
		panic(errors.Wrap(err, "open tree log file"))
	}
	flatFile, err = os.Create(filepath.Join(os.TempDir(), "errors_flat.log"))
	if err != nil {
		panic(errors.Wrap(err, "open flat log file"))
	}
	stdFile, err = os.Create(filepath.Join(os.TempDir(), "errors_std.log"))
	if err != nil {
		panic(errors.Wrap(err, "open std log file"))
	}
	txtCtxFile, err = os.Create(filepath.Join(os.TempDir(), "errors_txt.log"))
	if err != nil {
		panic(errors.Wrap(err, "open txt log file"))
	}

	treeLogger = slog.New(errorsctx.NewSLogHandlerTree(slog.NewJSONHandler(treeFile, &slog.HandlerOptions{})))
	flatLogger = slog.New(errorsctx.NewSLogHandlerFlat(slog.NewJSONHandler(flatFile, &slog.HandlerOptions{})))
	stdLogger = slog.New(slog.NewJSONHandler(stdFile, &slog.HandlerOptions{}))
	txtCtxLogger = slog.New(slog.NewJSONHandler(txtCtxFile, &slog.HandlerOptions{}))
}
