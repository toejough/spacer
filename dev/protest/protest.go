// Package protest provides procedure testing functionality.
package protest

import (
	"testing"
)

type FIFO[I any] struct {
	items []I
	name  string
}

func NewFIFO[I any](name string) *FIFO[I] {
	return &FIFO[I]{items: []I{}, name: name}
}

func (s *FIFO[I]) Len() int {
	return len(s.items)
}

func (s *FIFO[I]) Push(i I) {
	s.items = append(s.items, i)
}

func (s *FIFO[I]) MustPop(t *testing.T) I {
	t.Helper()

	if len(s.items) == 0 {
		t.Fatalf("expected to pop from %s stack, but there were no items in it\n", s.name)
	}

	var i I

	i, s.items = s.items[0], s.items[1:]

	return i
}

func RequireNext[I any](t *testing.T, expected I, fifo *FIFO[I], diff differ[I]) {
	t.Helper()

	if len(fifo.items) == 0 {
		t.Fatalf("expected to pop '%v' from '%s' stack, but there were no items in it\n", expected, fifo.name)
	}

	actual := fifo.MustPop(t)
	d := diff(expected, actual)

	if len(d) != 0 {
		t.Fatalf("expected next item in '%s' to be '%v', but a diff of '%s' was found\n", fifo.name, expected, d)
	}

	t.Logf("%s '%v'\n", fifo.name, actual)
}

type differ[T any] func(T, T) string

func RequireEmpty[I any](t *testing.T, s *FIFO[I]) {
	t.Helper()

	l := s.Len()
	if l != 0 {
		t.Fatalf("expected stack to be empty but it had %d items in it\n", l)
	}
}
