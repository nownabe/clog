package errors_test

import (
	"strings"
	"testing"

	"github.com/nownabe/clog/errors"
)

func TestNew(t *testing.T) {
	t.Parallel()

	err := errors.New("foo")

	var ews errors.ErrorWithStack
	if !errors.As(err, &ews) {
		t.Errorf("errors.As(err, &ews) should be true")
	}

	if ews.Stack() == nil {
		t.Errorf("ews.Stack() should not be nil")
	}

	stackLines := strings.Split(string(ews.Stack()), "\n")

	if len(stackLines) != 12 {
		t.Errorf("len(stackLines) got %d (%#v), want more than 12 lines", len(stackLines), stackLines)
	}

	if stackLines[0][:9] != "goroutine" {
		t.Errorf("stackLines[0][:9] got %q, want %q", stackLines[0][:9], "goroutine")
	}
}

func TestNewWithoutStack(t *testing.T) {
	t.Parallel()

	err := errors.NewWithoutStack("foo")

	var ews errors.ErrorWithStack
	if errors.As(err, &ews) {
		t.Errorf("errors.As(err, &ews) should be false")
	}
}

func TestErrorf(t *testing.T) {
	t.Parallel()

	err1 := errors.New("err1")
	err2 := errors.New("err1")
	err3 := errors.NewWithoutStack("err3")
	err4 := errors.NewWithoutStack("err4")
	wrap1 := func(f string, e error) error { return errors.Errorf(f, e) }
	wrap2 := func(f string, e ...any) error { return errors.Errorf(f, e...) }

	ews1 := err1.(errors.ErrorWithStack)

	assertEqual := func(t *testing.T, got, want []byte) {
		t.Helper()
		if len(got) != len(want) {
			t.Fatalf("len(got) got %d '%s', want %d '%s'", len(got), got, len(want), want)
		}
		for b := range got {
			if got[b] != want[b] {
				t.Errorf("got '%s', want '%s'", got, want)
			}
		}
	}

	tests := map[string]struct {
		err   error
		check func(t *testing.T, got []byte)
	}{
		"wrap error without stack": {
			err: wrap1("foo: %w", err3),
			check: func(t *testing.T, got []byte) {
				if !strings.Contains(string(got), "TestErrorf.func1") {
					t.Errorf("got '%s', should contain '%s'", got, "TestErrorf.func1")
				}
			},
		},
		"wrap error with stack": {
			err: wrap1("foo: %w", err1),
			check: func(t *testing.T, got []byte) {
				assertEqual(t, got, ews1.Stack())
			},
		},
		"wrap errors with stack and without stack": {
			err: wrap2("foo: %w, %w", err1, err3),
			check: func(t *testing.T, got []byte) {
				assertEqual(t, got, ews1.Stack())
			},
		},
		"wrap errors without stack and with stack": {
			err: wrap2("foo: %w, %w", err3, err1),
			check: func(t *testing.T, got []byte) {
				assertEqual(t, got, ews1.Stack())
			},
		},
		"wrap errors without stack": {
			err: wrap2("foo: %w, %w", err3, err4),
			check: func(t *testing.T, got []byte) {
				if !strings.Contains(string(got), "TestErrorf.func2") {
					t.Errorf("got '%s', should contain '%s'", got, "TestErrorf.func2")
				}
			},
		},
		"wrap errors with stack": {
			err: wrap2("foo: %w, %w", err1, err2),
			check: func(t *testing.T, got []byte) {
				assertEqual(t, got, ews1.Stack())
			},
		},
		"wrap no error": {
			err: wrap2("foo"),
			check: func(t *testing.T, got []byte) {
				if !strings.Contains(string(got), "TestErrorf.func2") {
					t.Errorf("got '%s', should contain '%s'", got, "TestErrorf.func2")
				}
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var ews errors.ErrorWithStack
			if !errors.As(tt.err, &ews) {
				t.Errorf("errors.As(tt.err, &ews) should be true")
			}

			tt.check(t, ews.Stack())
		})
	}
}

func TestErrorfWithoutStack(t *testing.T) {
	t.Parallel()

	err1 := errors.New("err1")
	err2 := errors.New("err1")
	err3 := errors.NewWithoutStack("err3")
	err4 := errors.NewWithoutStack("err4")
	wrap1 := func(f string, e error) error { return errors.ErrorfWithoutStack(f, e) }
	wrap2 := func(f string, e ...any) error { return errors.ErrorfWithoutStack(f, e...) }

	ews1 := err1.(errors.ErrorWithStack)

	assertEqual := func(t *testing.T, got, want []byte) {
		t.Helper()
		if len(got) != len(want) {
			t.Fatalf("len(got) got %d '%s', want %d '%s'", len(got), got, len(want), want)
		}
		for b := range got {
			if got[b] != want[b] {
				t.Errorf("got '%s', want '%s'", got, want)
			}
		}
	}

	tests := map[string]struct {
		err      error
		hasStack bool
		check    func(t *testing.T, got []byte)
	}{
		"wrap error without stack": {
			err:      wrap1("foo: %w", err3),
			hasStack: false,
		},
		"wrap error with stack": {
			err:      wrap1("foo: %w", err1),
			hasStack: true,
			check: func(t *testing.T, got []byte) {
				assertEqual(t, got, ews1.Stack())
			},
		},
		"wrap errors with stack and without stack": {
			err:      wrap2("foo: %w, %w", err1, err3),
			hasStack: true,
			check: func(t *testing.T, got []byte) {
				assertEqual(t, got, ews1.Stack())
			},
		},
		"wrap errors without stack and with stack": {
			err:      wrap2("foo: %w, %w", err3, err1),
			hasStack: true,
			check: func(t *testing.T, got []byte) {
				assertEqual(t, got, ews1.Stack())
			},
		},
		"wrap errors without stack": {
			err:      wrap2("foo: %w, %w", err3, err4),
			hasStack: false,
		},
		"wrap errors with stack": {
			err:      wrap2("foo: %w, %w", err1, err2),
			hasStack: true,
			check: func(t *testing.T, got []byte) {
				assertEqual(t, got, ews1.Stack())
			},
		},
		"wrap no error": {
			err:      wrap2("foo"),
			hasStack: false,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var ews errors.ErrorWithStack
			if errors.As(tt.err, &ews) {
				if !tt.hasStack {
					t.Fatal("errors.As(tt.err, &ews) should be false")
				}
				tt.check(t, ews.Stack())
			} else {
				if tt.hasStack {
					t.Error("errors.As(tt.err, &ews) should be true")
				}
			}
		})
	}
}
