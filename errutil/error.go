package errutil

import (
	"strings"

	"github.com/pkg/errors"
)

type customError struct {
	originErr      error
	originStackErr error
	currentErr     error
}

func newCustomError(err error) *customError {
	if err == nil {
		return nil
	}
	stackErr := errors.Wrap(err, "Error")
	return &customError{
		originErr:      err,
		originStackErr: stackErr,
		currentErr:     stackErr,
	}
}

func (inst *customError) Error() string {
	// return inst.originStackErr.Error()
	// return SerrorStackTrace(inst)
	return inst.currentErr.Error()
}

func (inst *customError) Wrap(message string) error {
	inst.currentErr = errors.Wrap(inst.currentErr, message)
	return inst
}

func (inst *customError) Wrapf(message string, data ...interface{}) error {
	inst.currentErr = errors.Wrapf(inst.currentErr, message, data...)
	return inst
}

func New(err error) error {
	return newCustomError(err)
}

func NewWithMessage(messageFormat string, data ...interface{}) error {
	return newCustomError(errors.Errorf(messageFormat, data...))
}

func IsError(err error, originalError error) bool {
	if customError, ok := err.(*customError); ok {
		return errors.Cause(customError.currentErr) == errors.Cause(originalError)
	}
	return err == originalError
}

func Wrap(err error, message string) error {
	if customError, ok := err.(*customError); ok {
		return customError.Wrap(message)
	}
	return newCustomError(err).Wrap(message)
}

func Wrapf(err error, message string, data ...interface{}) error {
	if customError, ok := err.(*customError); ok {
		return customError.Wrapf(message, data...)
	}
	return newCustomError(err).Wrapf(message, data...)
}

func Combine(errs []error) error {
	errStrs := []string{}
	for _, err := range errs {
		if err != nil {
			errStrs = append(errStrs, err.Error())
		}
	}
	if len(errStrs) == 0 {
		return nil
	}
	return NewWithMessage(strings.Join(errStrs, ",\n"))
}
