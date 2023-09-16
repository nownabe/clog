/*
Package errors provides errors with stack trace.
*/
package errors

import (
	"errors"
	"fmt"
	"runtime"
)

const stackBufSize = 1024

// ErrorWithStack is an interface that has Stack method to provide stack trace of en error.
// The format of stack should be same as runtime.Stack.
// See https://pkg.go.dev/runtime#Stack.
type ErrorWithStack interface {
	Stack() []byte
	Error() string
}

// New returns an error with stack.
func New(text string) error {
	return newWithStack(errors.New(text))
}

// NewWithoutStack returns an error without stack.
func NewWithoutStack(text string) error {
	return errors.New(text)
}

// As is the same as errors.As.
// See https://pkg.go.dev/errors#As.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Is is the same as errors.Is.
// See https://pkg.go.dev/errors#Is.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Join is the same as errors.Join.
// See https://pkg.go.dev/errors#Join.
func Join(errs ...error) error {
	return errors.Join(errs...)
}

// Unwrap is the same as errors.Unwrap.
// See https://pkg.go.dev/errors#Unwrap.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Errorf returns an error with stack.
// %w can be used as well as fmt.Errorf.
// See https://pkg.go.dev/fmt#Errorf.
func Errorf(format string, args ...any) error {
	if hasStack(args...) {
		return fmt.Errorf(format, args...)
	}

	return newWithStack(fmt.Errorf(format, args...))
}

// ErrorfWithoutStack returns an error without stack.
// %w can be used as well as fmt.Errorf.
// See https://pkg.go.dev/fmt#Errorf.
func ErrorfWithoutStack(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

// WithStack just wrap the given error with stack.
func WithStack(err error) error {
	if err == nil {
		return nil
	}

	return newWithStack(err)
}

type withStack struct {
	err   error
	stack []byte
}

func newWithStack(err error) *withStack {
	buf := make([]byte, stackBufSize)
	n := runtime.Stack(buf, false)
	return &withStack{
		err:   err,
		stack: buf[:n],
	}
}

func (e *withStack) Error() string {
	return e.err.Error()
}

func (e *withStack) Unwrap() error {
	return e.err
}

func (e *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(s, "%v\n%s", e.err, e.stack)
	case 's':
		fmt.Fprint(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *withStack) Stack() []byte {
	return e.stack
}

func hasStack(args ...any) bool {
	for _, a := range args {
		if err, ok := a.(error); ok {
			var ews ErrorWithStack
			if errors.As(err, &ews) {
				return true
			}
		}
	}

	return false
}
