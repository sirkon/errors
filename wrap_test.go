package errors_test

import (
	"fmt"
	"testing"

	"github.com/sirkon/errors"
)

func ExampleWrap() {
	fmt.Println(errors.Wrap(errors.Const("example"), "error"))

	// Output:
	// error: example
}

func ExampleWrapf() {
	fmt.Println(errors.Wrapf(errors.Const("example"), "formatted error"))

	// Output:
	// formatted error: example
}

func TestWrapPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			return
		}

		t.Errorf("wrap on nil error must cause a panic")
	}()

	t.Log(errors.Wrap(nil, "error"))
}

func TestWrapfPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			return
		}

		t.Errorf("wrapf on nil error must cause a panic")
	}()

	t.Log(errors.Wrapf(nil, "error"))
}
