package errors_test

import (
	"fmt"

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
