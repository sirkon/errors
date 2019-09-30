package errors

import (
	"bytes"
)

var _ error = List{}

// List list of errors
type List []error

func (l List) Error() string {
	switch len(l) {
	case 0:
		// специально вызываем панику – так положено, ибо нужно явно проверять
		var err error
		return err.Error()
	case 1:
		return l[0].Error()
	default:
		var buf bytes.Buffer
		for i, err := range l {
			if i > 0 {
				buf.WriteString("; ")
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// As applies As function with given target to each error in a list until success
func (l List) As(target interface{}) bool {
	for _, err := range l {
		if As(err, target) {
			return true
		}
	}
	return false
}

// Is just like As Is applies Is function with given err to each error in a list until success
func (l List) Is(err error) bool {
	for _, e := range l {
		if Is(e, err) {
			return true
		}
	}
	return false
}

// And creates a new combined error. Avoid `err = errors.And(err1, err)` and use `err = errors.And(err, err1)` instead
// when collecting numerous errors. This will save some memory allocations.
func And(err1 error, err2 error) error {
	if err1 == nil {
		return err2
	}
	if err2 == nil {
		return err1
	}
	switch v := err1.(type) {
	case List:
		return append(v, err2)
	default:
		res := make(List, 2, 4)
		res[0] = err1
		res[1] = err2
		return res
	}
}
