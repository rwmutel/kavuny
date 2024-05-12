package main

type HttpError struct {
	err        error
	msg        string
	returnCode int
}

func (e *HttpError) Error() string {
	return e.msg + ": " + e.err.Error()
}

func (e *HttpError) Unwrap() error {
	return e.err
}

func (e *HttpError) Code() int {
	return e.returnCode
}

func NewHttpError(err error, msg string, returnCode int) *HttpError {
	return &HttpError{err: err, msg: msg, returnCode: returnCode}
}
