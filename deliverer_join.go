//go:build sirkon_errors_hidden

package errors

// TODO move this functionality out of the package as it is rather specific for the deliverer user

// NewJoinDeliverer constructs JoinDeliverer
func NewJoinDeliverer(src ErrorContextDeliverer) *JoinDeliverer {
	return &JoinDeliverer{src: src}
}

// JoinDeliverer a composite deliverer joining values with same name into a slice of values
type JoinDeliverer struct {
	src ErrorContextDeliverer
}

// Deliver to implement deliverer
func (d *JoinDeliverer) Deliver(cons ErrorContextConsumer) {
	c := &joiningConsumer{
		order:  nil,
		single: map[string]interface{}{},
		mult:   map[string][]interface{}{},
	}

	d.src.Deliver(c)

	for _, name := range c.order {
		v, ok := c.single[name]
		if !ok {
			cons.Any(name, c.mult[name])
			continue
		}

		switch vv := v.(type) {
		case bool:
			cons.Bool(name, vv)
		case int:
			cons.Int(name, vv)
		case int8:
			cons.Int8(name, vv)
		case int16:
			cons.Int16(name, vv)
		case int32:
			cons.Int32(name, vv)
		case int64:
			cons.Int64(name, vv)
		case uint:
			cons.Uint(name, vv)
		case uint8:
			cons.Uint8(name, vv)
		case uint16:
			cons.Uint16(name, vv)
		case uint32:
			cons.Uint32(name, vv)
		case uint64:
			cons.Uint64(name, vv)
		case float32:
			cons.Float32(name, vv)
		case float64:
			cons.Float64(name, vv)
		case string:
			cons.String(name, vv)
		default:
			cons.Any(name, v)
		}
	}
}

func (d *JoinDeliverer) Error() string {
	return ""
}

var (
	_ ErrorContextDeliverer = &JoinDeliverer{}
)

type joiningConsumer struct {
	order  []string
	single map[string]interface{}
	mult   map[string][]interface{}
}

func (j *joiningConsumer) val(name string, value interface{}) {
	if v, ok := j.single[name]; ok {
		j.mult[name] = []interface{}{v, value}
		delete(j.single, name)
		return
	}

	if v, ok := j.mult[name]; ok {
		j.mult[name] = append(v, value)
		return
	}

	j.order = append(j.order, name)
	j.single[name] = value
}

func (j *joiningConsumer) Bool(name string, value bool) {
	j.val(name, value)
}

func (j *joiningConsumer) Int(name string, value int) {
	j.val(name, value)
}

func (j *joiningConsumer) Int8(name string, value int8) {
	j.val(name, value)
}

func (j *joiningConsumer) Int16(name string, value int16) {
	j.val(name, value)
}

func (j *joiningConsumer) Int32(name string, value int32) {
	j.val(name, value)
}

func (j *joiningConsumer) Int64(name string, value int64) {
	j.val(name, value)
}

func (j *joiningConsumer) Uint(name string, value uint) {
	j.val(name, value)
}

func (j *joiningConsumer) Uint8(name string, value uint8) {
	j.val(name, value)
}

func (j *joiningConsumer) Uint16(name string, value uint16) {
	j.val(name, value)
}

func (j *joiningConsumer) Uint32(name string, value uint32) {
	j.val(name, value)
}

func (j *joiningConsumer) Uint64(name string, value uint64) {
	j.val(name, value)
}

func (j *joiningConsumer) Float32(name string, value float32) {
	j.val(name, value)
}

func (j *joiningConsumer) Float64(name string, value float64) {
	j.val(name, value)
}

func (j *joiningConsumer) String(name string, value string) {
	j.val(name, value)
}

func (j *joiningConsumer) Any(name string, value interface{}) {
	j.val(name, value)
}
