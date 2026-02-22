package errors_test

import (
	"fmt"

	"github.com/sirkon/errors"
	"github.com/sirkon/errors/errorsctx"
)

func ExampleGetContextDeliverer() {
	err := errors.New("test error").Int("count", 12).Str("text", "Hello")
	err = errors.Wrap(err, "check").Uint("unsign", 333).Bytes("actually-text", []byte("World!"))
	err = errors.Just(err).Bool("this-is-the-last", true).Bytes("real-raw", []byte{1, 2, 3})

	dlvr := errors.GetContextDeliverer(err)
	var c errorsctx.Consumer
	dlvr.Deliver(&c)

	for _, layer := range c.Layers {
		fmt.Println(layer)
		for _, pair := range layer.Pairs {
			fmt.Printf("    %s: %v\n", pair.Key, pair.Value.Any())
		}
	}

	// Output:
	// NEW: test error
	//     count: 12
	//     text: Hello
	// WRAP: check
	//     unsign: 333
	//     actually-text: World!
	// CTX
	//     this-is-the-last: true
	//     real-raw: [1 2 3]
}
