// Package protest provides procedure testing functionality.
package protest

import (
	"fmt"
	"testing"
)

type FIFO[I any] struct {
	items []I
}

func (s *FIFO[I]) Len() int {
	return len(s.items)
}

func (s *FIFO[I]) Push(i I) {
	s.items = append(s.items, i)
}

var ErrPop = fmt.Errorf("unable to pop from stack: no items in it")

func (s *FIFO[I]) Pop() (I, error) {
	if len(s.items) == 0 {
		zeroValue := *new(I) //nolint:gocritic // cannot do I(nil) with generics

		return zeroValue, ErrPop
	}

	var i I

	i, s.items = s.items[0], s.items[1:]

	return i, nil
}

func (s *FIFO[I]) MustPop(t *testing.T) I {
	t.Helper()

	if len(s.items) == 0 {
		t.Fatalf("unable to pop from stack: no items in it")
	}

	var i I

	i, s.items = s.items[0], s.items[1:]

	return i
}

func stringDiff(e, a string) string {
	if e != a {
		return fmt.Sprintf("%s != %s", e, a)
	}

	return ""
}

func RequireCall(t *testing.T, expected, actual string) {
	t.Helper()

	require(t, expected, actual, stringDiff, "call")
}


type differ[T any] func(T, T) string

func RequireArgs[T any](t *testing.T, expected, actual T, d differ[T]) {
	t.Helper()

	require(t, expected, actual, d, "args")
}

func RequireReturn[T any](t *testing.T, expectedArgs, actualArgs T, d differ[T]) {
	t.Helper()

	require(t, expectedArgs, actualArgs, d, "return")
}

func require[T any](t *testing.T, expected, actual T, d differ[T], what string) {
	t.Helper()

	diff := d(expected, actual)
	if len(diff) != 0 {
		t.Fatalf("expected %s to be '%v', but a diff of '%s' was found\n", what, expected, diff)
	}

	t.Logf("%s '%v'\n", what, actual)
}

func RequireEmpty[I any](t *testing.T, s FIFO[I]) {
	t.Helper()

	l := s.Len()
	if l != 0 {
		t.Fatalf("expected stack to be empty but it had %d items in it\n", l)
	}
}
