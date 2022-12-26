// Package protest provides procedure testing functionality.
package protest

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type FIFO[I any] struct {
	items chan I
	name  string
}

func NewFIFO[I any](name string) *FIFO[I] {
	return &FIFO[I]{items: make(chan I), name: name}
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

var ErrTimedOut = fmt.Errorf("timed out")

func (s *FIFO[I]) WaitForNext(d time.Duration) (next I, err error) {
	select {
	case next = <-s.items:
		return next, nil
	case <-time.After(d):
		return next, fmt.Errorf("waited %v for an item from %s FIFO, but there was none: %w", d, s.name, ErrTimedOut)
	}
}

func (s *FIFO[I]) MustWaitForNext(t *testing.T, d time.Duration) (next I) {
	t.Helper()

	var err error
	next, err = s.WaitForNext(1 * time.Second)

	if err != nil {
		t.Fatal(err)
	}

	return next
}

func (s *FIFO[I]) GetNext() (next I, err error) {
	return s.WaitForNext(1 * time.Second)
}

func (s *FIFO[I]) MustGetNext(t *testing.T) (next I) {
	t.Helper()

	return s.MustWaitForNext(t, 1*time.Second)
}

func (s *FIFO[I]) RequireNext(t *testing.T, expected I) {
	t.Helper()

	next := s.MustGetNext(t)
	if !reflect.DeepEqual(expected, next) {
		t.Fatalf("expected %#v but found %#v instead", expected, next)
	}
}

func (s *FIFO[I]) RequireNextWithin(t *testing.T, expected I, d time.Duration) {
	t.Helper()

	next := s.MustWaitForNext(t, d)
	if reflect.DeepEqual(expected, next) {
		t.Fatalf("expected %#v but found %#v instead", expected, next)
	}
}
