// Package protest provides procedure testing functionality.
package protest

import (
	"testing"
	"time"
)

type FIFODeps[I any] struct {
	Differ differ[I]
	T      *testing.T
}

type FIFO[I any] struct {
	items chan I
	name  string
	deps  FIFODeps[I]
}

func NewFIFO[I any](name string, deps FIFODeps[I]) *FIFO[I] {
	return &FIFO[I]{items: make(chan I), name: name, deps: deps}
}

func (s *FIFO[I]) Close() {
    close(s.items)
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

func (s *FIFO[I]) GetNext() I {
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
