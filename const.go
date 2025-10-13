package errors

// Const is a type implementing the error interface and capable of being a constant.
type Const string

func (c Const) Error() string {
	return string(c)
}

// Is for [errors.Is].
func (c Const) Is(err error) bool {
	v, ok := err.(Const)
	if !ok {
		return false
	}

	return v == c
}

// As for [errors.As].
func (c Const) As(target any) bool {
	v, ok := target.(*Const)
	if !ok {
		return false
	}

	*v = c
	return true
}
