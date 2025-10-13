package errors_test

import (
	"io"
	"testing"

	"github.com/sirkon/errors"
)

func ExampleJust() {
	err := errors.Just(io.EOF).
		Str("name", "value").
		Int("value", 1)
	LogError(err)

	// Output:
	// error message: EOF
	//  - name: value
	//  - value: 1
}

func TestJustPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			return
		}

		t.Error("just on nil error must cause a panic")
	}()

	t.Log(errors.Just(nil))
}
