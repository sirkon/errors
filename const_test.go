package errors_test

import (
	"fmt"

	"github.com/sirkon/errors"
)

func ExampleConst_Is() {
	err := errors.Wrap(errors.Const("error"), "message")

	fmt.Println(errors.Is(err, errors.Const("error")))
	fmt.Println(errors.Is(err, errors.Const("another error")))
	fmt.Println(errors.Is(errors.New("error"), errors.Const("error")))

	// Output:
	// true
	// false
	// false
}
