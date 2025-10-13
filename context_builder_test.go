package errors_test

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/sirkon/errors"
)

func ExampleCtx() {
	// Basic demo to show how context is stored and extracted.

	errctx := errors.Ctx().
		Bool("true", true).
		Int("int", 12).
		Int8("i8", math.MinInt8).
		Int16("i16", math.MinInt16).
		Int32("i32", math.MinInt32).
		Int64("i64", math.MinInt64).
		Uint("u", math.MaxUint).
		Uint8("u8", math.MaxUint8).
		Uint16("u16", math.MaxUint16).
		Uint32("u32", math.MaxUint32).
		Uint64("u64", math.MaxUint64).
		Flt32("e", math.E).
		Flt64("pi", math.Pi).
		Str("str", "hello world").
		Stg("stringer", testStringer{}).
		Strs("strlist", strings.Split("1 2 3", " ")).
		Bytes("bytes-printable", []byte("123")).
		Bytes("bytes-raw", []byte{1, 2, 3}).
		Type("type", new(bytes.Buffer)).
		Any("map", map[int]int{1: 3, 2: 2, 3: 1})

	err := errors.New("error").
		Int("code", 404).
		Pfx("ctx").
		WithCtx(errctx)

	LogError(err)

	// Output:
	// error message: error
	//  - code: 404
	//  - ctx-true: true
	//  - ctx-int: 12
	//  - ctx-i8: -128
	//  - ctx-i16: -32768
	//  - ctx-i32: -2147483648
	//  - ctx-i64: -9223372036854775808
	//  - ctx-u: 18446744073709551615
	//  - ctx-u8: 255
	//  - ctx-u16: 65535
	//  - ctx-u32: 4294967295
	//  - ctx-u64: 18446744073709551615
	//  - ctx-e: 2.7182817
	//  - ctx-pi: 3.141592653589793
	//  - ctx-str: hello world
	//  - ctx-stringer: test stringer
	//  - ctx-strlist: [1 2 3]
	//  - ctx-bytes-printable: 123
	//  - ctx-bytes-raw: [1 2 3]
	//  - ctx-type: *bytes.Buffer
	//  - ctx-map: map[1:3 2:2 3:1]
}

func ExampleContext_WithCtx() {
	// Why simply using a variable for context won't work.
	ctxBase := errors.Ctx().Int("code", 404)
	ctxMsg := ctxBase.Str("msg", "not found")
	ctxBool := ctxBase.Bool("found", false)
	LogError(errors.New("error-code").WithCtx(ctxBase))
	LogError(errors.New("error-msg").WithCtx(ctxMsg))
	LogError(errors.New("error-bool").WithCtx(ctxBool))

	fmt.Println()
	fmt.Println("...reuse context properly...")
	ctxBase = errors.Ctx().Int("code", 403)
	ctxMsg = errors.CtxFrom(ctxBase).Str("msg", "forbidden")
	ctxBool = errors.CtxFrom(ctxBase).Bool("allowed", false)
	LogError(errors.New("error-code").WithCtx(ctxBase))
	LogError(errors.New("error-msg").WithCtx(ctxMsg))
	LogError(errors.New("error-bool").WithCtx(ctxBool))

	// Output:
	// error message: error-code
	//  - code: 404
	//  - msg: not found
	//  - found: false
	// error message: error-msg
	//  - code: 404
	//  - msg: not found
	//  - found: false
	// error message: error-bool
	//  - code: 404
	//  - msg: not found
	//  - found: false
	//
	// ...reuse context properly...
	// error message: error-code
	//  - code: 403
	// error message: error-msg
	//  - code: 403
	//  - msg: forbidden
	// error message: error-bool
	//  - code: 403
	//  - allowed: false
}
