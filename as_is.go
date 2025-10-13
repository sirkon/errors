package errors

import "errors"

// Is is a wrapper for [errors.Is] from the standard library.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As is a wrapper for [errors.As] from the standard library.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// AsType is a type-safe generic version of [As], allowing usage without pre-declaring a variable.
func AsType[T error](err error) (T, bool) {
	var t T
	if !errors.As(err, &t) {
		return t, false
	}

	return t, true
}
