package errors

// Just is for when annotation is not needed, but structured context is.
// This function solves that problem by returning an error with the text from the given one
// and allowing to add structured values.
func Just(err error) *Error {
	if err == nil {
		_ = err.Error()
	}

	res := &Error{
		msg:       "",
		err:       err,
		ctxPrefix: "",
		ctx:       nil,
	}

	if insertLocations {
		res.setLoc()
	}

	return res
}
