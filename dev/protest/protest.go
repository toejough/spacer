// Package protest provides procedure testing functionality.
package protest

import (
	"testing"
	"time"
)

type NewFIFODeps[I any] struct {
	Differ differ[I]
	T      *testing.T
}

type FIFO[I any] struct {
	items chan I
	name  string
	deps  NewFIFODeps[I]
}

func NewFIFO[I any](name string) *FIFO[I] {
	return &FIFO[I]{items: make(chan I), name: name}
}

func NewFIFO2[I any](name string, deps NewFIFODeps[I]) *FIFO[I] {
	return &FIFO[I]{items: make(chan I), name: name, deps: deps}
}

func (s *FIFO[I]) Len() int {
	return len(s.items)
}

func (s *FIFO[I]) Push(i I) {
	s.items <- i
}

func (s *FIFO[I]) RequireNext(next I) {
	s.deps.T.Helper()

	select {
	case i := <-s.items:
		d := s.deps.Differ(next, i)

		if len(d) != 0 {
			s.deps.T.Fatalf("expected next item in '%s' to be '%v', but a diff of '%s' was found\n", s.name, next, d)
		}

		s.deps.T.Logf("%s '%v'\n", s.name, i)
	case <-time.After(1 * time.Second):
		s.deps.T.Fatalf("expected to pop from %s FIFO, but there were no items in it after 1s of waiting.\n", s.name)
		panic("panic here to satisfy linter")
	}
}

func (s *FIFO[I]) RequireEmpty() {
	RequireEmpty(s.deps.T, s)
}

func (s *FIFO[I]) MustPop2() I {
	s.deps.T.Helper()

	select {
	case i := <-s.items:
		return i
	case <-time.After(1 * time.Second):
		s.deps.T.Fatalf("expected to pop from %s FIFO, but there were no items in it after 1s of waiting.\n", s.name)
		panic("panic here to satisfy linter")
	}
}


type differ[T any] func(T, T) string

func RequireEmpty[I any](t *testing.T, s *FIFO[I]) {
	t.Helper()

	l := s.Len()
	if l != 0 {
		t.Fatalf("expected stack to be empty but it had %d items in it\n", l)
	}
}
