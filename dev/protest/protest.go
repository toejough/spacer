// Package protest provides procedure testing functionality.
package protest

import (
	"fmt"
	"reflect"
	"time"
)

type FIFODeps[I any] struct {
	Differ differ[I]
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

var ErrChannelNotClosed = fmt.Errorf("the channel isn't closed")

func (s *FIFO[I]) RequireClosedAndEmpty() error {
	select {
	case value := <-s.items:
		if !reflect.ValueOf(value).IsZero() {
			return fmt.Errorf("expected no more values in %s, but found %v: %w", s.name, value, ErrChannelNotClosed)
		}
	case <-time.After(1 * time.Second):
		return fmt.Errorf("expected %s to be closed, but it was not after 1s of waiting: %w", s.name, ErrChannelNotClosed)
	}

	return nil
}

type differ[T any] func(T, T) string

var ErrTimedOut = fmt.Errorf("timed out")

func (s *FIFO[I]) WaitForNext(d time.Duration) (next I, err error) {
	select {
	case next = <-s.items:
		return next, nil
	case <-time.After(d):
		return next, fmt.Errorf("waited %v for an item from %s FIFO, but there was none: %w", d, s.name, ErrTimedOut)
	}
}

func (s *FIFO[I]) GetNext() (next I, err error) {
	return s.WaitForNext(1 * time.Second)
}
