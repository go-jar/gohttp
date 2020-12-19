package query

import (
	"github.com/go-jar/goerror"
)

type Value interface {
	Required() bool
	Set(str string) error
	Check() bool
	Error() *goerror.Error
}

type baseValue struct {
	required bool
	errno    string
	msg      string
}

func newBaseValue(required bool, errno string, msg string) *baseValue {
	return &baseValue{
		required: required,
		errno:    errno,
		msg:      msg,
	}
}

func (b *baseValue) Required() bool {
	return b.required
}

func (b *baseValue) Error() *goerror.Error {
	return goerror.New(b.errno, b.msg)
}
