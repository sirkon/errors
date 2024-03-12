package errors

import "errors"

// Join simply wraps stdlib errors.Join.
func Join(errs ...error) error {
	return errors.Join(errs...)
}
