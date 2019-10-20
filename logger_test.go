package errors_test

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/sirkon/errors"
)

func ExampleLogger() {
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
	err := errors.Context().Int64("epoch", 1234).String("user-id", "1").New("no such user this time")
	err = errors.Wrap(err, "check if user can be synced")
	errors.Logger(logger, err).Error().Msg("syncing user data failed")
	err = errors.Context().
		Int("i", 1).
		Int8("i8", 1).
		Int16("i16", 1).
		Int32("i32", 1).
		Int64("i64", 1).
		Uint("u", 1).
		Uint8("u8", 1).
		Uint16("u16", 1).
		Uint32("u32", 1).
		Uint64("u64", 1).
		String("s", "string").
		Strings("ss", []string{"1"}).
		Any("a", map[string]int{}).
		Line().
		New("message")
	errors.Logger(logger, err).Error().Msg("all types")

	// Output:
	// {"level":"error","error":"check if user can be synced: no such user this time","epoch":1234,"user-id":"1","message":"syncing user data failed"}
	// {"level":"error","error":"message","err-pos":"/home/emacs/go/src/github.com/sirkon/errors/logger_test.go:30","i":1,"i8":1,"i16":1,"i32":1,"i64":1,"u":1,"u8":1,"u16":1,"u32":1,"u64":1,"s":"string","ss":["1"],"a":{},"message":"all types"}
}
