package errors

import "errors"

// Join is a wrapper for [errors.Join] from the standard library.
func Join(errs ...error) error {
	return errors.Join(errs...)
}
