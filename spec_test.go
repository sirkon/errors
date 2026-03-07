package errors_test

import (
	"fmt"
	"io"

	"github.com/sirkon/errors"
)

func ExampleSpec() {
	var err error
	err = errors.Wrap(io.EOF, "inner check 1")
	err = errors.Spec(err, new(1))
	err = fmt.Errorf("foreign check: %w", err)
	err = errors.Wrap(err, "inner check 2")
	err = errors.Wrap(err, "inner check 3")
	fmt.Println(err, asSpec(err), errors.IsSpec[*int](err))

	err = fmt.Errorf("foreign check: %w", io.EOF)
	err = errors.Spec(err, new(2))
	fmt.Println(err, asSpec(err), errors.IsSpec[*int](err))

	err = fmt.Errorf("foreign check: %w", io.EOF)
	fmt.Println(err, asSpec(err), errors.IsSpec[*int](err))

	// Output:
	// inner check 3: inner check 2: foreign check: inner check 1: EOF 1 true
	// foreign check: EOF 2 true
	// foreign check: EOF -1 false
}

func asSpec(err error) int {
	res, ok := errors.AsSpec[*int](err)
	if !ok {
		return -1
	}

	return *res
}
