package errors_test

import (
	stderrors "errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/sirkon/errors"
)

var count int

func BenchmarkErrorsWrapFixed(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		err := getErrorsWrapErrorNoContext()
		count += len(err.Error())
	}
}

func BenchmarkFmtErrorfFixed(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		err := getFmtErrorfFixed()
		count += len(err.Error())
	}
}

func BenchmarkErrorsWrapf(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		err := getErrorsWrapf()
		count += len(err.Error())
	}
}

func BenchmarkErrorsWrapContext(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		err := getErrorsWrapContext()
		count += len(err.Error())
	}
}

func BenchmarkFmtErrorfShortContext(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		err := getFmtErrorfShortContext()
		count += len(err.Error())
	}
}

func BenchmarkErrorsWrapLongContext(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		err := getErrorsWrapLargerContext()
		count += len(err.Error())
	}
}

func BenchmarkFmtErrorfLongContext(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		err := getFmtErrorfLargerContext()
		count += len(err.Error())
	}
}

func getErrorsWrapErrorNoContext() error {
	var err error = errors.New("some error")
	err = errors.Wrap(err, "wrap 1")
	err = errors.Wrap(err, "wrap 2")
	err = errors.Wrap(err, "wrap 3")
	err = errors.Wrap(err, "wrap 4")

	return err
}

func getFmtErrorfFixed() error {
	err := stderrors.New("some error")
	err = fmt.Errorf("wrap 1: %w", err)
	err = fmt.Errorf("wrap 2: %w", err)
	err = fmt.Errorf("wrap 3: %w", err)
	err = fmt.Errorf("wrap 4: %w", err)

	return err
}

func getErrorsWrapf() error {
	err := os.ErrNotExist
	filename := "query.sql"
	dbname := "db-1"
	domain := "sales"
	deparmentID := "deparment-123"

	err = errors.Wrapf(err, "open query file %s", filename)
	err = errors.Wrapf(err, "query users in DB %s", dbname)
	err = errors.Wrapf(err, "get domain %s users", domain)
	err = errors.Wrapf(err, "create iterator over deparment %s users", deparmentID)

	return err
}

func getFmtErrorfShortContext() error {
	err := os.ErrNotExist
	filename := "query.sql"
	dbname := "db-1"
	domain := "sales"
	deparmentID := "deparment-123"

	err = fmt.Errorf("open query file %s: %w", filename, err)
	err = fmt.Errorf("query users in DB %s: %w", dbname, err)
	err = fmt.Errorf("get domain %s users: %w", domain, err)
	err = fmt.Errorf("create iterator over deparment %s users: %w", deparmentID, err)

	return err
}

func getErrorsWrapContext() error {
	err := errors.New("some error").Int("some-int", 12).Str("some-str", "hello")
	err = errors.Wrap(err, "wrap 1").Str("wrap-str", "world")
	err = errors.Wrap(err, "wrap 2")
	err = errors.Wrap(err, "wrap 3").Str("method", "Method")

	return err
}

func getErrorsWrapLargerContext() error {
	err := errors.Wrap(io.ErrClosedPipe, "read stream").Int("buf-length", 4096)
	err = errors.Wrap(err, "read with retries").
		F64("average-failure-response-delay", 0.12).
		Int("how-many-retries", 4)
	err = errors.Wrapf(err, "read GetFile stream")
	err = errors.Wrap(err, "transfer data").
		Int("transfered-before-failure", 16384).
		Int("retries-made", 9)
	err = errors.Wrap(err, "get image data").Str("data-id", "6142a749-aaa2-4383-b6bd-9d0adfd9d330")
	err = errors.Wrap(err, "resize image").
		Str("img-type", "image/png").
		F64("scale", 1.2).
		Str("sharpening-type", "GAUSSIAN")
	err = errors.Wrap(err, "process user avatar")
	err = errors.Wrapf(err, "finish %s", "user-create-routine").
		Str("user-id", "e0e3804f-b0f7-4fc6-a995-fd20c4994810")
	err = errors.Wrap(err, "replay wal session").
		Str("wal-session-id", "0bf25c6a-d5d6-4b08-b381-9c6a26ea55c0").
		Int("session-replay-no", 3).
		F64("replay-duration", 0.87323).
		Int("replays-left", 7)

	return err
}

func getFmtErrorfLargerContext() error {
	err := fmt.Errorf("read %d bytes of stream: %w", 4096, io.ErrClosedPipe)
	err = fmt.Errorf("read with retries average-delay[%f] retries-done[%d]: %w", 0.12, 4, err)
	err = fmt.Errorf("read GetFile stream: %w", err)
	err = fmt.Errorf("transfer data bytes-completed[%d] total-retries[%d]: %w", 16384, 9, err)
	err = fmt.Errorf("get image data from record %q: %w", "6142a749-aaa2-4383-b6bd-9d0adfd9d330", err)
	err = fmt.Errorf("resize %s image with the scale %f: %w", "image/jpeg", 1.2, err)
	err = fmt.Errorf("process user avatar: %w", err)
	err = fmt.Errorf(
		"finish process[%s] user-id[%s]: %w",
		"user-create-routine",
		"e0e3804f-b0f7-4fc6-a995-fd20c4994810",
		err,
	)
	err = fmt.Errorf(
		"replay wal session=%q session-replay-no=%d replay-duration=%f replays-left=%d: %w",
		"0bf25c6a-d5d6-4b08-b381-9c6a26ea55c0",
		3,
		0.87323,
		7,
		err,
	)

	return err
}

func init() {
	errors.DoNotInsertLocations()
}
