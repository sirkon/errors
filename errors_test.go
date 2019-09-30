package errors_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/sirkon/errors"
)

func panicRequired(t *testing.T, action func()) (panicHappened bool) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("panic was expected")
		}
	}()
	action()
	return
}

func ExampleWrap() {
	fmt.Println(errors.Wrap(io.EOF, "error"))
	fmt.Println(errors.Wrap(errors.Wrapf(io.EOF, "msg %d", 1), "error"))

	// Output:
	// error: EOF
	// error: msg 1: EOF
}

func ExampleWrappedError_Unwrap() {
	err := errors.Wrap(io.EOF, "error 1")
	err = errors.Wrap(err, "error 2")
	err = errors.Wrap(err, "error 3")

	fmt.Println(err)
	fmt.Println(err.(interface{ Unwrap() error }).Unwrap())

	// Output:
	// error 3: error 2: error 1: EOF
	// EOF
}

func ExampleAs() {
	var ce customError
	err := errors.Wrap(customError("message"), "message")
	err = customWrapper{err: err}
	if errors.As(err, &ce) {
		fmt.Println(ce)
	}
	if !errors.As(err, &fakeError{}) {
		fmt.Println(err)
	}
	errs := errors.And(io.EOF, err)
	errs = errors.And(errs, io.EOF)
	if errors.As(errs, &ce) {
		fmt.Println(ce)
	}

	fmt.Println(errors.As(errs, &fakeError{}))

	// Output:
	// error message
	// custom wrap: message: error message
	// error message
	// false
}

func ExampleIs() {
	ce := customError("message")
	err := errors.Wrap(customError("message"), "message")
	err = customWrapper{err: err}

	fmt.Println(errors.Is(err, ce))
	fmt.Println(errors.Is(err, fakeError{}))

	errs := errors.And(io.EOF, err)
	errs = errors.And(errs, io.EOF)
	fmt.Println(errors.Is(errs, ce))
	fmt.Println(errors.Is(errs, fakeError{}))

	fmt.Println("checking nil cases")
	fmt.Println(errors.Is(nil, nil))
	fmt.Println(errors.Is(io.EOF, nil))
	fmt.Println(errors.Is(nil, io.EOF))

	// Output:
	// true
	// false
	// true
	// false
	// checking nil cases
	// true
	// false
	// false
}

func ExampleAnd() {
	fmt.Println(errors.And(io.EOF, io.EOF))
	fmt.Println(errors.And(nil, io.EOF))
	fmt.Println(errors.And(io.EOF, nil))
	fmt.Println(errors.List{io.EOF})

	// Output:
	// EOF; EOF
	// EOF
	// EOF
	// EOF
}

func ExampleNew() {
	fmt.Println(errors.New("error"))

	// Output:
	// error
}

func ExampleNewf() {
	fmt.Println(errors.Newf("error %s", "message"))

	// Output:
	// error message
}

// panic raising

func TestAs_WithPanic(t *testing.T) {
	panicRequired(t, func() {
		errors.As(io.EOF, nil)
	})
	panicRequired(t, func() {
		errors.As(io.EOF, 1)
	})
	panicRequired(t, func() {
		var tmp interface{}
		ttt := 1
		tmp = &ttt
		errors.As(io.EOF, tmp)
	})
}

func TestWrap_WithPanic(t *testing.T) {
	action := func() {
		t.Log(errors.Wrap(nil, "error"))
	}
	panicRequired(t, action)
}

func TestList_WithPanic(t *testing.T) {
	var l errors.List
	panicRequired(t, func() {
		fmt.Println(l.Error())
	})
}

// testing stuff

type customWrapper struct {
	err error
}

func (cw customWrapper) Error() string { return "custom wrap: " + cw.err.Error() }
func (cw customWrapper) Unwrap() error { return cw.err }

type customError string

func (ce customError) Error() string { return "error " + string(ce) }

type fakeError struct{}

func (ce fakeError) Error() string { return "fake error" }
