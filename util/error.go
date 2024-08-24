package util

import (
	"fmt"
)

type CustomError struct {
	code string
	msg  string
}

func NewCustomError(code, msg string) CustomError {
	return CustomError{
		code: code,
		msg:  msg,
	}
}

func (e CustomError) Error() string {
	return e.msg
}

func (e CustomError) Code() string {
	return e.code
}

func (e CustomError) ErrorWithCode() string {
	return fmt.Sprintf("%s:%s", e.code, e.msg)
}

func (e CustomError) Is(target error) bool {
	t, ok := target.(CustomError)
	return ok && t.code == e.code
}
