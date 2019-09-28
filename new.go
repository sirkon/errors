package errors

import (
	"errors"
	"fmt"
)

// New just an alias for stdlib's New
func New(msg string) error {
	return errors.New(msg)
}

// Newf an alias for fmt.Errorf
func Newf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}
