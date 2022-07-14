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

func RequireCall(t *testing.T, expectedCall string, actualCall string) {
	t.Helper()

	if actualCall != expectedCall {
		t.Fatalf("expected call to '%s', but '%s' was called instead\n", expectedCall, actualCall)
	}

	t.Logf("called '%s'\n", actualCall)
}

func RequireEmpty[I any](t *testing.T, s FIFO[I]) {
	t.Helper()

	l := s.Len()
	if l != 0 {
		t.Fatalf("expected stack to be empty but it had %d items in it\n", l)
	}
}

type Diffable[D any] interface {
	Diff(D) string
}

func RequireArgs[D Diffable[D]](t *testing.T, expectedArgs D, actualArgs D) {
	t.Helper()

	diff := expectedArgs.Diff(actualArgs)
	if len(diff) != 0 {
		t.Fatalf("expected args to be '%v', but a diff of '%s' was found\n", expectedArgs, diff)
	}

	t.Logf("args '%v'\n", actualArgs)
}
