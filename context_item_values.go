package errors

type (
	intValue     int
	int8Value    int8
	int16Value   int16
	int32Value   int32
	int64Value   int64
	uintValue    uint
	uint8Value   uint8
	uint16Value  uint16
	uint32Value  uint32
	uint64Value  uint64
	stringValue  string
	float32Value float32
	float64Value float64
	anyValue     struct {
		val interface{}
	}
	stringsValue []string
)

func (i intValue) reportItem(name string, builder ContextReportBuilder) {
	builder.Int(name, int(i))
}

func (i int8Value) reportItem(name string, builder ContextReportBuilder) {
	builder.Int8(name, int8(i))
}

func (i int16Value) reportItem(name string, builder ContextReportBuilder) {
	builder.Int16(name, int16(i))
}

func (i int32Value) reportItem(name string, builder ContextReportBuilder) {
	builder.Int32(name, int32(i))
}

func (i int64Value) reportItem(name string, builder ContextReportBuilder) {
	builder.Int64(name, int64(i))
}

func (i uintValue) reportItem(name string, builder ContextReportBuilder) {
	builder.Uint(name, uint(i))
}

func (i uint8Value) reportItem(name string, builder ContextReportBuilder) {
	builder.Uint8(name, uint8(i))
}

func (i uint16Value) reportItem(name string, builder ContextReportBuilder) {
	builder.Uint16(name, uint16(i))
}

func (i uint32Value) reportItem(name string, builder ContextReportBuilder) {
	builder.Uint32(name, uint32(i))
}

func (i uint64Value) reportItem(name string, builder ContextReportBuilder) {
	builder.Uint64(name, uint64(i))
}

func (s stringValue) reportItem(name string, builder ContextReportBuilder) {
	builder.String(name, string(s))
}

func (i float32Value) reportItem(name string, builder ContextReportBuilder) {
	builder.Float32(name, float32(i))
}

func (i float64Value) reportItem(name string, builder ContextReportBuilder) {
	builder.Float64(name, float64(i))
}

func (a anyValue) reportItem(name string, builder ContextReportBuilder) {
	builder.Any(name, a.val)
}

func (ss stringsValue) reportItem(name string, builder ContextReportBuilder) {
	builder.Strings(name, ss)
}

//

var (
	_ itemValue = intValue(0)
	_ itemValue = int8Value(0)
	_ itemValue = int16Value(0)
	_ itemValue = int32Value(0)
	_ itemValue = int64Value(0)
	_ itemValue = uint8Value(0)
	_ itemValue = uintValue(0)
	_ itemValue = uint16Value(0)
	_ itemValue = uint32Value(0)
	_ itemValue = uint64Value(0)
	_ itemValue = stringValue("")
	_ itemValue = float32Value(0)
	_ itemValue = float64Value(0)
	_ itemValue = anyValue{}
	_ itemValue = stringsValue{}
)
