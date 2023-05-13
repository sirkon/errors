package errors

// Just a simple wrapper over an error that allows to add
// structured context values without adding any text info.
func Just(err error) Error {
	if err == nil {
		err.Error()
	}

	v, ok := err.(Error)
	if ok {
		return v
	}

	return Error{
		msg:       "",
		err:       err,
		ctxPrefix: "",
		ctx:       nil,
	}
}
