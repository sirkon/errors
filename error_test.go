package errors_test

import (
	"fmt"
	"io"

	"awesome-errors"
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
	}
	for _, err := range errs {
		fmt.Println(err)
	}

	// Output:
	// wrap no 2: wrap: new error
	// wrap: EOF
}
