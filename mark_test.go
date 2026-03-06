package errors_test

import (
	"fmt"
	"io"

	"github.com/sirkon/errors"
)

type tagType struct{}

func ExampleMark() {
	var err error
	err = errors.Wrap(io.EOF, "inner check 1")
	err = errors.Mark(err, &tagType{})
	err = fmt.Errorf("foreign check: %w", err)
	err = errors.Wrap(err, "inner check 2")
	err = errors.Wrap(err, "inner check 3")
	fmt.Println(err, hasTagMark(err))

	err = fmt.Errorf("foreign check: %w", io.EOF)
	err = errors.Mark(err, &tagType{})
	fmt.Println(err, hasTagMark(err))

	err = fmt.Errorf("foreign check: %w", io.EOF)
	fmt.Println(err, hasTagMark(err))

	// Output:
	// inner check 3: inner check 2: foreign check: inner check 1: EOF true
	// foreign check: EOF true
	// foreign check: EOF false
}

func hasTagMark(err error) bool {
	_, ok := errors.HasMark[*tagType](err)
	return ok
}
