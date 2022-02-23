// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import "errors"

// Is errors.Is wrapper
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As errors.As wrapper
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Dig looks for error of type T in err's chain. Returns a point to an error found
// or nil otherwise.
func Dig[T error](err error) *T {
	var t T
	if !errors.As(err, &t) {
		return nil
	}

	return &t
}
