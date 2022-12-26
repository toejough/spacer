// Package protest provides procedure testing functionality.
//
// FIFO pop flavors:
// (Must)Pop(As|EqualTo)(Within)
// All pops: pop the next item off of the FIFO
// Must option:
// * With: pass a *testing.T as the first arg, and failure will cause a t.Fatal.
// * Without: failures will be returned as an error
// As/EqualTo options:
// * As: pass a pointer to a type, and the next item will be assigned to that pointer if possible
// * EqualTo: pass a value, and the next item will be checked for deep-equality
// * Without: the next item will be returned
// Within option:
// * With: pass a duration as the final arg, and a timeout will only occur after that duration
// * Without: timeouts will occur after 1s
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

var (
	ErrChannelNotClosed     = fmt.Errorf("the channel isn't closed")
	ErrNotEqual             = fmt.Errorf("values not equal")
	ErrTimedOut             = fmt.Errorf("timed out")
	ErrNilPointerTarget     = fmt.Errorf("target must be a non-nil pointer")
	ErrUnassignableToTarget = fmt.Errorf("value is unassignable to target")
)

// TODO Drop the fifo name.
// NewFIFO creates a new FIFO and returns a pointer to it.
func NewFIFO[I any](name string) *FIFO[I] {
	return &FIFO[I]{items: make(chan I), name: name}
}

// Close closes the underlying resources for the FIFO.
func (s *FIFO[I]) Close() {
	close(s.items)
}

// Push pushes the given value into the FIFO.
func (s *FIFO[I]) Push(i I) {
	s.items <- i
}

// Pop pops the next thing from the FIFO, waiting up to 1s for it to be available. If it is available within the
// timeout, it returns the value and a nil error. If it is not available within the timeout, it returns ErrTimedOut.
func (s *FIFO[I]) Pop() (next I, err error) {
	return s.PopWithin(1 * time.Second)
}

// PopEqualTo pops the next thing from the FIFO, waiting up to 1s for it to be available. If it is available within the
// timeout and equal to the expected value, it returns a nil error. Equality is tested with reflect.DeepEqual. If it is
// not available within the timeout, it returns ErrTimedOut. If it is not equal, it returns ErrNotEqual.
func (s *FIFO[I]) PopEqualTo(expected I) (err error) {
	return s.PopEqualToWithin(expected, 1*time.Second)
}

// PopAs pops the next thing from the FIFO, waiting up to 1s for it to be available. If it is available within the
// timeout, it attempts to set the target to be the value. If settable, it returns a nil error. If it is not available
// within the timeout, it returns ErrTimedOut. If it is not settable, it returns ErrUnassignableToTarget.
func (s *FIFO[I]) PopAs(target I) (err error) {
	return s.PopAsWithin(target, 1*time.Second)
}

// PopWithin pops the next thing from the FIFO, waiting up to the given duration for it to be available. If it is
// available, it returns the value and a nil error. If it is not available within the timeout, it returns ErrTimedOut.
func (s *FIFO[I]) PopWithin(d time.Duration) (next I, err error) {
	select {
	case next = <-s.items:
		return next, nil
	case <-time.After(d):
		return next, fmt.Errorf("waited %v for an item from %s FIFO, but there was none: %w", d, s.name, ErrTimedOut)
	}
}

// PopEqualToWithin pops the next thing from the FIFO, waiting up to the given duration for it to be available. If it is
// available and equal to the expected value, it returns a nil error. Equality is tested with reflect.DeepEqual. If it
// is not available within the timeout, it returns ErrTimedOut. If it is not equal, it returns ErrNotEqual.
func (s *FIFO[I]) PopEqualToWithin(expected I, d time.Duration) (err error) {
	next, err := s.PopWithin(d)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, next) {
		return fmt.Errorf("expected %#v but found %#v instead: %w", expected, next, ErrNotEqual)
	}

	return nil
}

// PopAsWithin pops the next thing from the FIFO, waiting up to the given duration for it to be available. If it is
// available, it attempts to set the target to be the value. If settable, it returns a nil error. If it is not
// available within the timeout, it returns ErrTimedOut. If it is not settable, it returns ErrUnassignableToTarget.
func (s *FIFO[I]) PopAsWithin(target I, d time.Duration) (err error) {
	next, err := s.PopWithin(d)
	if err != nil {
		return err
	}

	// most of this copied from the errors.As implementation at
	// https://cs.opensource.google/go/go/+/refs/tags/go1.19.4:src/errors/wrap.go;l=78
	val := reflect.ValueOf(target)
	typ := val.Type()

	if typ.Kind() != reflect.Ptr || val.IsNil() {
		return ErrNilPointerTarget
	}

	targetType := typ.Elem()

	if reflect.TypeOf(next).AssignableTo(targetType) {
		val.Elem().Set(reflect.ValueOf(next))
		return nil
	}

	return fmt.Errorf("%#v can not be set as %T: %w", next, target, ErrUnassignableToTarget)
}

// MustPop pops the next thing from the FIFO, waiting up to 1s for it to be available. If it is available within the
// timeout, it returns the value. If it is not available within the timeout, it triggers a fatal test failure.
func (s *FIFO[I]) MustPop(t *testing.T) (next I) {
	t.Helper()

	return s.MustPopWithin(t, 1*time.Second)
}

// MustPopEqualTo pops the next thing from the FIFO, waiting up to 1s for it to be available. If it is available within
// the timeout and equal to the expected value, it returns. Equality is tested with reflect.DeepEqual. If it is not
// available within the timeout, it triggers a fatal test failure. If it is not equal, it triggers a fatal test failure.
func (s *FIFO[I]) MustPopEqualTo(t *testing.T, expected I) {
	t.Helper()

	s.MustPopEqualToWithin(t, expected, 1*time.Second)
}

// MustPopAs pops the next thing from the FIFO, waiting up to 1s for it to be available. If it is available within the
// timeout, it attempts to set the target to be the value. If it is not available within the timeout, it triggers a
// fatal test failure. If it is not settable, it triggers a fatal test failure.
func (s *FIFO[I]) MustPopAs(t *testing.T, target I) {
	t.Helper()

	s.MustPopAsWithin(t, target, 1*time.Second)
}

// MustPopWithin pops the next thing from the FIFO, waiting up to the given duration for it to be available. If it is
// available, it returns the value. If it is not available within the timeout, it triggers a fatal test failure.
func (s *FIFO[I]) MustPopWithin(t *testing.T, d time.Duration) (next I) {
	t.Helper()

	var err error

	next, err = s.PopWithin(d)

	if err != nil {
		t.Fatal(err)
	}

	return
}

// MustPopEqualToWithin pops the next thing from the FIFO, waiting up to the given duration for it to be available. If
// it is available and equal to the expected value, it returns. Equality is tested with reflect.DeepEqual. If it is not
// available within the timeout, it triggers a fatal test failure. If it is not equal, it triggers a fatal test failure.
func (s *FIFO[I]) MustPopEqualToWithin(t *testing.T, expected I, d time.Duration) {
	t.Helper()

	err := s.PopEqualToWithin(expected, d)
	if err != nil {
		t.Fatal(err)
	}
}

// MustPopAsWithin pops the next thing from the FIFO, waiting up to the given duration for it to be available. If it is
// available, it attempts to set the target to be the value. If settable, it returns. If it is not available within the
// timeout, it triggers a fatal test failure. If it is not settable, it triggers a fatal test failure.
func (s *FIFO[I]) MustPopAsWithin(t *testing.T, target I, d time.Duration) {
	t.Helper()

	err := s.PopAsWithin(target, d)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *FIFO[I]) RequireClosedAndEmpty() error {
	select {
	// TODO use value, ok := ... "ok" will be false if the channel is closed and empty
	case value := <-s.items:
		if reflect.ValueOf(value).IsValid() && !reflect.ValueOf(value).IsZero() {
			return fmt.Errorf("expected no more values in %s, but found %v: %w", s.name, value, ErrChannelNotClosed)
		}
	case <-time.After(1 * time.Second):
		return fmt.Errorf("expected %s to be closed, but it was not after 1s of waiting: %w", s.name, ErrChannelNotClosed)
	}

	return nil
}

// TODO RequireDeepEqual
