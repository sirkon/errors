package errors_test

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"github.com/sirkon/errors"
)

func ExampleNew() {
	fmt.Println(errors.New("error example"))
	// output:
	// error example
}

func ExampleNewf() {
	fmt.Println(errors.Newf("error %s", "example"))
	// output:
	// error example
}

func ExampleJust() {
	err := errors.Just(io.EOF).Str("name", "value")
	d := errors.GetContextDeliverer(err)
	var cons testContextConsumer
	d.Deliver(cons)
	fmt.Println(err.Error() == io.EOF.Error())

	err = errors.Just(errors.New("error").Int("value", 1))
	d = errors.GetContextDeliverer(err)
	cons = testContextConsumer{}
	d.Deliver(cons)
	fmt.Println(err.Error())

	// Output:
	// name: value
	// true
	// value: 1
	// error
}

func ExampleWrap() {
	fmt.Println(errors.Wrap(errors.Const("example"), "error"))
	// output:
	// error: example
}

func ExampleWrapf() {
	fmt.Println(errors.Wrapf(errors.Const("example"), "formatted error"))
	// output:
	// formatted error: example
}

func ExampleDig() {
	err := fmt.Errorf("example: %w", errors.New("error"))
	if pe := errors.Dig[errors.Error](err); pe != nil {
		fmt.Println((*pe).Error())
	}

	err = fmt.Errorf("example: %w", errors.Const("const error"))
	if pe := errors.Dig[errors.Error](err); pe != nil {
		fmt.Println((*pe).Error())
	}
	// output:
	// error
}

func ExampleIs() {
	err := errors.Wrap(io.EOF, "read file")
	if !errors.Is(err, io.EOF) {
		fmt.Println("must not be here")
	}

	e := errors.New("i am an error")
	if !errors.Is(errors.Wrap(e, "covered"), e) {
		fmt.Println("must not be here again")
	}
	// output:
}

func ExampleGetContextDeliverer() {
	err := errors.New("error").
		Bool("b", true).
		Int("i", -1).
		Int8("i8", math.MaxInt8).
		Int16("i16", math.MinInt16).
		Int32("i32", math.MinInt32).
		Int64("i64", math.MinInt64).
		Uint("u", 1).
		Uint8("u8", math.MaxUint8).
		Uint16("u16", math.MaxUint16).
		Uint32("u32", math.MaxUint32).
		Uint64("u64", math.MaxUint64).
		Float32("f32", math.MaxFloat32).
		Float64("f64", math.MaxFloat64).
		Str("string", "str").
		Strs("strings", []string{"1", "2", "3"}).
		Stg("stringer", testStringer{}).
		Any("object", map[string]int{
			"key": 12,
		}).
		Any("bytes", []byte("123")).
		Type("type-name", bytes.NewBuffer(nil))
	err = errors.Wrap(err, "wrapping context").Str("wrapped", "value")
	d := errors.GetContextDeliverer(err)
	var cons testContextConsumer
	d.Deliver(cons)

	if errors.GetContextDeliverer(errors.Const("error")) != nil {
		fmt.Println("must no be here")
	}
	// output:
	// wrapped: value
	// b: true
	// i: -1
	// i8: 127
	// i16: -32768
	// i32: -2147483648
	// i64: -9223372036854775808
	// u: 1
	// u8: 255
	// u16: 65535
	// u32: 4294967295
	// u64: 18446744073709551615
	// f32: 3.4028235e+38
	// f64: 1.7976931348623157e+308
	// string: str
	// strings: [1 2 3]
	// stringer: test stringer
	// object: map[key:12]
	// bytes: 123
	// type-name: *bytes.Buffer
}

type testContextConsumer struct{}

func (t testContextConsumer) Bool(name string, value bool) {
	fmt.Printf("%s: %v\n", name, value)
}
func (t testContextConsumer) Int(name string, value int) {
	fmt.Printf("%s: %v\n", name, value)
}
func (t testContextConsumer) Int8(name string, value int8) {
	fmt.Printf("%s: %v\n", name, value)
}
func (t testContextConsumer) Int16(name string, value int16) {
	fmt.Printf("%s: %v\n", name, value)
}
func (t testContextConsumer) Int32(name string, value int32) {
	fmt.Printf("%s: %v\n", name, value)
}
func (t testContextConsumer) Int64(name string, value int64) {
	fmt.Printf("%s: %v\n", name, value)
}
func (t testContextConsumer) Uint(name string, value uint) {
	fmt.Printf("%s: %v\n", name, value)
}
func (t testContextConsumer) Uint8(name string, value uint8) {
	fmt.Printf("%s: %v\n", name, value)
}
func (t testContextConsumer) Uint16(name string, value uint16) {
	fmt.Printf("%s: %v\n", name, value)
}
func (t testContextConsumer) Uint32(name string, value uint32) {
	fmt.Printf("%s: %v\n", name, value)
}
func (t testContextConsumer) Uint64(name string, value uint64) {
	fmt.Printf("%s: %v\n", name, value)
}
func (t testContextConsumer) Float32(name string, value float32) {
	fmt.Printf("%s: %v\n", name, value)
}
func (t testContextConsumer) Float64(name string, value float64) {
	fmt.Printf("%s: %v\n", name, value)
}
func (t testContextConsumer) String(name string, value string) {
	fmt.Printf("%s: %v\n", name, value)
}
func (t testContextConsumer) Any(name string, value interface{}) {
	fmt.Printf("%s: %v\n", name, value)
}

type testStringer struct{}

func (testStringer) String() string {
	return "test stringer"
}
