package check

import "fmt"

type validateError struct {
	msg       string
	permanent bool
}

func newValidateError(format string, a ...interface{}) *validateError {
	return &validateError{
		msg: fmt.Sprintf(format, a...),
	}
}

func (err *validateError) AsPermanent() *validateError {
	err.permanent = true
	return err
}
