package check

import "fmt"

type validateError struct {
	msg              string
	rateLimitReached bool
}

func newValidateError(format string, a ...interface{}) *validateError {
	return &validateError{
		msg: fmt.Sprintf(format, a...),
	}
}

func (err *validateError) RateLimitReached() *validateError {
	err.rateLimitReached = true
	return err
}
