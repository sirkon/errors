package errors_test

import (
	"fmt"
	"io"

	"github.com/sirkon/errors"
)

func ExampleError_Error() {
	errs := []error{
		func() error {
			err := errors.New("new error")
			err = errors.Wrap(err, "wrap")
			err = errors.Wrapf(err, "wrap no %d", 2)
			return err
		}(),
		func() error {
			return errors.Wrap(io.EOF, "wrap")
		}(),
		func() error {
			var err error
			err = errors.Wrap(io.EOF, "wrap")
			err = fmt.Errorf("foreign wrap: %w", err)
			err = errors.Spec(err, new(0))
			return err
		}(),
	}
	for _, err := range errs {
		fmt.Println(err)
	}

	// Output:
	// wrap no 2: wrap: new error
	// wrap: EOF
	// foreign wrap: wrap: EOF
}
