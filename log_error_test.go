package errors_test

import (
	"fmt"
	"go/token"

	"github.com/sirkon/errors"
)

// LogError prints the error text and its context.
func LogError(err error) {
	fmt.Println("error message:", err.Error())

	deliverer := errors.GetContextDeliverer(err)
	if deliverer == nil {
		return
	}

	var cons testConsumer
	deliverer.Deliver(&cons)
	for _, val := range cons.ctx {
		fmt.Println(" - "+val.name+":", val.value)
	}
}

type contextTuple struct {
	name  string
	value any
}

type testConsumer struct {
	loc []token.Position
	ctx []contextTuple
}

func (t *testConsumer) NextLink() {}

func (t *testConsumer) SetLinkInfo(token.Position, errors.ErrorChainLinkDescriptor) {}

func (t *testConsumer) Bool(name string, value bool) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) Int(name string, value int) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) Int8(name string, value int8) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) Int16(name string, value int16) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) Int32(name string, value int32) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) Int64(name string, value int64) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) Uint(name string, value uint) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) Uint8(name string, value uint8) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) Uint16(name string, value uint16) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) Uint32(name string, value uint32) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) Uint64(name string, value uint64) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) Flt32(name string, value float32) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) Flt64(name string, value float64) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) Str(name string, value string) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) Any(name string, value any) {
	t.ctx = append(t.ctx, contextTuple{
		name:  name,
		value: value,
	})
}

func (t *testConsumer) AddLoc(loc token.Position, _ errors.ErrorChainLinkDescriptor) {
	t.loc = append(t.loc, loc)
}
