package errors_test

import (
	"bytes"
	"fmt"
	"math"
	"strconv"

	"github.com/sirkon/errors"
)

func ExampleGetContextDeliverer() {
	// Handle the case of an error without structured context.
	if errors.GetContextDeliverer(errors.Const("error")) != nil {
		fmt.Println("must not be here")
	}

	// Now the main case where we need to show values from the context.
	// Here we add the prefix "error" – call Pfx("error") – and add
	// a bunch of values of various types.
	err := errors.New("error").Pfx("").Pfx("context").
		Bool("b", true).
		Int("i", -1).
		Int8("i8", math.MinInt8).
		Int16("i16", math.MinInt16).
		Int32("i32", math.MinInt32).
		Int64("i64", math.MinInt64).
		Uint("u", 1).
		Uint8("u8", math.MaxUint8).
		Uint16("u16", math.MaxUint16).
		Uint32("u32", math.MaxUint32).
		Uint64("u64", math.MaxUint64).
		Flt32("f32", math.MaxFloat32).
		Flt64("f64", math.MaxFloat64).
		Str("string", "str").
		Strs("strings", []string{"1", "2", "3"}).
		Stg("stringer", testStringer{}).
		Any("object", map[string]int{
			"key": 12,
		}).
		Bytes("bytes", []byte("123")).Bytes("bytes-raw", []byte{1, 2, 3}).
		Type("type-name", bytes.NewBuffer(nil))

	// Add a bit more context, this time without a prefix.
	err = errors.Wrap(err, "wrapping context").Str("wrapped", "value")
	// Print the error text.
	fmt.Println("error message:", strconv.Quote(err.Error()))

	// Get the context deliverer and feed it our ready implementation
	// of the context receiver, then output the accumulated data to STDOUT with a bit of formatting.
	cons := &testConsumer{}
	d := errors.GetContextDeliverer(err)
	d.Deliver(cons)
	for _, ctx := range cons.ctx {
		fmt.Println(" - "+ctx.name+":", ctx.value)
	}

	// output:
	// error message: "wrapping context: error"
	//  - wrapped: value
	//  - context-b: true
	//  - context-i: -1
	//  - context-i8: -128
	//  - context-i16: -32768
	//  - context-i32: -2147483648
	//  - context-i64: -9223372036854775808
	//  - context-u: 1
	//  - context-u8: 255
	//  - context-u16: 65535
	//  - context-u32: 4294967295
	//  - context-u64: 18446744073709551615
	//  - context-f32: 3.4028235e+38
	//  - context-f64: 1.7976931348623157e+308
	//  - context-string: str
	//  - context-strings: [1 2 3]
	//  - context-stringer: test stringer
	//  - context-object: map[key:12]
	//  - context-bytes: 123
	//  - context-bytes-raw: [1 2 3]
	//  - context-type-name: *bytes.Buffer
}

type testStringer struct{}

func (testStringer) String() string {
	return "test stringer"
}

func init() {
	// Enable to slightly improve coverage.
	errors.InsertLocations()
}
