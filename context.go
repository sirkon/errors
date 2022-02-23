package errors

// TODO move this out of this package

// ErrorContextConsumer an abstraction meant to consume structured context of an error.
type ErrorContextConsumer interface {
	Bool(name string, value bool)
	Int(name string, value int)
	Int8(name string, value int8)
	Int16(name string, value int16)
	Int32(name string, value int32)
	Int64(name string, value int64)
	Uint(name string, value uint)
	Uint8(name string, value uint8)
	Uint16(name string, value uint16)
	Uint32(name string, value uint32)
	Uint64(name string, value uint64)
	Float32(name string, value float32)
	Float64(name string, value float64)
	String(name string, value string)
	Any(name string, value interface{})
}

// ErrorContextDeliverer an abstraction to work with ErrorContextConsumer, it delivers
// structured context variables into a consumer.
type ErrorContextDeliverer interface {
	Deliver(cons ErrorContextConsumer)
	Error() string
}
