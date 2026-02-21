package errors

import (
	"errors"
)

// Is mirrors [errors.Is].
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// As mirrors [errors.As].
func As(err error, target any) bool {
	return errors.As(err, target)
}

// AsType mirrors [errors.AsType].
func AsType[E error](err error) (E, bool) {
	return errors.AsType[E](err)
}

// Join mirrors [errors.Join].
func Join(errs ...error) error {
	return errors.Join(errs...)
}
