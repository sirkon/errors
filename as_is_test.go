package errors_test

import (
	"fmt"
	"io"

	"github.com/sirkon/errors"
)

func ExampleIs() {
	err := errors.Wrap(io.ErrNoProgress, "read data")
	fmt.Println(errors.Is(err, io.ErrNoProgress))
	fmt.Println(errors.Is(
		errors.Wrap(errors.New("error"), "wrapped"),
		io.EOF,
	))

	// Output:
	// true
	// false
}

func ExampleAs() {
	err := fmt.Errorf("read data: %w", errors.New("root error"))

	var e *errors.Error
	if errors.As(err, &e) {
		fmt.Println(e.Error())
	}

	var ae errors.Const
	if errors.As(err, &ae) {
		fmt.Println(e.Error())
	}

	// Output:
	// root error
}

func ExampleAsType() {
	err := fmt.Errorf("read data: %w", errors.New("root error"))

	if v, ok := errors.AsType[*errors.Error](err); ok {
		fmt.Println(v)
	}

	if _, ok := errors.AsType[errors.Const](err); ok {
		fmt.Println("must not be here with", err)
	}

	// output:
	// root error
}
