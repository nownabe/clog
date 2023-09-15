package errors

import (
	"errors"
	"fmt"
	"runtime"
)

const stackBufSize = 1024

type ErrorWithStack interface {
	Stack() []byte
}

func New(text string) error {
	return newWithStack(errors.New(text))
}

func NewWithoutStack(text string) error {
	return errors.New(text)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Errorf(format string, args ...any) error {
	if hasStack(args...) {
		return fmt.Errorf(format, args...)
	}

	return newWithStack(fmt.Errorf(format, args...))
}

func ErrorfWithoutStack(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

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
