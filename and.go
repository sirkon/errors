package errors

import (
	"strings"
)

var _ error = List{}

// List list of errors
type List []error

func (ce List) Error() string {
	var buf strings.Builder
	for i, err := range ce {
		if i > 0 {
			buf.WriteString("; ")
		}
		buf.WriteString(err.Error())
	}
	return buf.String()
}

// As applies As function with given target to each error in a list until success
func (ce List) As(target interface{}) bool {
	for _, err := range ce {
		if As(err, target) {
			return true
		}
	}
	return false
}

// Is just like As Is applies Is function with given err to each error in a list until success
func (ce List) Is(err error) bool {
	for _, e := range ce {
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
