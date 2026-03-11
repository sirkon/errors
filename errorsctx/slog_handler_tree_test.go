package errorsctx_test

import (
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sirkon/errors"
	"github.com/sirkon/errors/errorsctx"
)

var (
	treeLogger    *slog.Logger
	flatLogger    *slog.Logger
	stdLogger     *slog.Logger
	txtCtxLogger  *slog.Logger
	discardLogger *slog.Logger

	treeFile       *os.File
	flatFile       *os.File
	stdFile        *os.File
	txtCtxFile     *os.File
	benchWriteFile *os.File
)

func TestMain(t *testing.M) {
	var files []*os.File
	defer func() {
		for _, file := range files {
			if err := file.Close(); err != nil {
				fmt.Println(errors.Wrap(err, "close "+file.Name()))
			}
		}
	}()
	treeFile = createFile("errors_tree.log")
	files = append(files, treeFile)
	flatFile = createFile("errors_flat.log")
	files = append(files, flatFile)
	stdFile = createFile("errors_std.log")
	files = append(files, stdFile)
	txtCtxFile = createFile("errors_txt.log")
	files = append(files, txtCtxFile)
	benchWriteFile = createFile("bench_write.log")
	files = append(files, benchWriteFile)

	treeLogger = slog.New(errorsctx.NewSLogHandlerTree(slog.NewJSONHandler(treeFile, &slog.HandlerOptions{})))
	flatLogger = slog.New(errorsctx.NewSLogHandlerFlat(slog.NewJSONHandler(flatFile, &slog.HandlerOptions{})))
	stdLogger = slog.New(slog.NewJSONHandler(stdFile, &slog.HandlerOptions{}))
	txtCtxLogger = slog.New(slog.NewJSONHandler(txtCtxFile, &slog.HandlerOptions{}))
	discardLogger = slog.New(errorsctx.NewSLogHandlerFlat(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{})))

	t.Run()
}

func BenchmarkErrorsTree(b *testing.B) {
	b.ReportAllocs()
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
	b.ReportAllocs()
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
	b.ReportAllocs()
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
	b.ReportAllocs()
	for b.Loop() {
		err := fmt.Errorf("this is an error bytes[%v] text-bytes[%s]", []byte{1, 2, 3}, "Hello World!")
		err = fmt.Errorf("check error count[%d] is-wrap-layer[%v]: %w", 333, true, err)
		err = fmt.Errorf("context pi[%g] e[%g]: %w", math.Pi, math.E, err)

		txtCtxLogger.Error("failed to do something 1", slog.Any("err", err))
	}
}

func BenchmarkWriteCost(b *testing.B) {
	b.ReportAllocs()
	line := strings.Repeat("1", 1024)
	for b.Loop() {
		if _, err := benchWriteFile.WriteString(line); err != nil {
			b.Error("failed to write file:", err)
			return
		}
	}
}

func BenchmarkAssembleAndFormattingCost(b *testing.B) {
	b.ReportAllocs()
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

		discardLogger.Error("failed to do something", slog.Any("err", err))
	}
}

func BenchmarkAssembleCost(b *testing.B) {
	b.ReportAllocs()
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

		if err == nil {
			panic("must not be nil")
		}
	}
}

func BenchmarkKeyCost(b *testing.B) {
	b.ReportAllocs()
	dst := make([]byte, 1024)
	key := "6142a749-aaa2-4383-b6bd-9d0adfd9d330"

	for b.Loop() {
		dst = dst[:0]

		//
		dst = binary.AppendUvarint(dst, uint64(len(key)))
		dst = append(dst, key...)
	}
}

func BenchmarkAttrCost(b *testing.B) {
	b.ReportAllocs()
	dst := make([]byte, 1024)
	key := "6142a749-aaa2-4383-b6bd-9d0adfd9d330"

	var i int
	for b.Loop() {
		dst = dst[:0]

		dst = binary.AppendUvarint(dst, uint64(len(key)))
		dst = append(dst, key...)
		dst = binary.LittleEndian.AppendUint64(dst, uint64(i))
	}
}

func createFile(name string) *os.File {
	file, err := os.Create(filepath.Join(os.TempDir(), name))
	if err != nil {
		panic(errors.Wrap(err, "create file "+name))
	}

	return file
}
