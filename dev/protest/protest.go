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
	"time"
)

type FIFO[I any] struct {
	items   chan I
	name    string
	oneShot bool
}

var (
	ErrChannelNotClosed     = fmt.Errorf("the channel isn't closed")
	ErrFIFOClosed           = fmt.Errorf("the FIFO is closed")
	ErrChannelNotEmpty      = fmt.Errorf("the channel isn't empty")
	ErrNilPointerTarget     = fmt.Errorf("target must be a non-nil pointer")
	ErrNotEqual             = fmt.Errorf("values not equal")
	ErrTimedOut             = fmt.Errorf("timed out")
	ErrUnassignableToTarget = fmt.Errorf("value is unassignable to target")
)

// NewFIFO creates a new FIFO and returns a pointer to it.
func NewFIFO[I any](name string) *FIFO[I] {
	return &FIFO[I]{items: make(chan I), name: name, oneShot: false}
}

// NewOneShotFIFO creates a new FIFO and returns a pointer to it. The FIFO is a "one-shot", meaning it only allows one
// push, and one pop. After the first push, subsequent pushes will panic.
func NewOneShotFIFO[I any](name string) *FIFO[I] {
	return &FIFO[I]{items: make(chan I), name: name, oneShot: true}
}

// Close closes the underlying resources for the FIFO.
func (s *FIFO[I]) Close() {
	close(s.items)
}

// Push pushes the given value into the FIFO.
func (s *FIFO[I]) Push(i I) {
	s.items <- i
	if s.oneShot {
		s.Close()
	}
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
func (s *FIFO[I]) PopWithin(duration time.Duration) (next I, err error) {
	var open bool
	select {
	case next, open = <-s.items:
		if !open {
			return next, fmt.Errorf("could not pop an item from %s: %w", s.name, ErrFIFOClosed)
		}

		if s.oneShot {
			err = s.ConfirmClosedWithin(duration)
		}

		return next, err
	case <-time.After(duration):
		return next, fmt.Errorf("waited %v for an item from %s FIFO, but there was none: %w", duration, s.name, ErrTimedOut)
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

	return Equal(expected, next)
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

type Tester interface {
	Helper()
	Fatal(...any)
}

// MustPop pops the next thing from the FIFO, waiting up to 1s for it to be available. If it is available within the
// timeout, it returns the value. If it is not available within the timeout, it triggers a fatal test failure.
func (s *FIFO[I]) MustPop(t Tester) (next I) {
	t.Helper()

	return s.MustPopWithin(t, 1*time.Second)
}

// MustPopNamed pops the next thing from the FIFO, waiting up to 1s for it to be available. If it is available within
// the timeout, and it's an AnyCall with the expected name, it returns the value. If it is not available within the
// timeout, it triggers a fatal test failure. If it is not an AnyCall, it triggers a fatal test error. If it is not the
// expected name, it triggers a fatal test error.
func (s *FIFO[I]) MustPopNamed(test Tester, name string) (next I) {
	test.Helper()

	call, err := s.PopWithin(1 * time.Second)
	if err != nil {
		test.Fatal(fmt.Errorf("didn't find the expected call to %s: %w", name, err))
	}

	anyCall, ok := any(call).(AnyCall)
	if !ok {
		test.Fatal(fmt.Sprintf("didn't find the expected call to %s: value popped in MustPopNamed was not an AnyCall", name))
	}

	MustEqual(test, name, anyCall.Name)

	return call
}

// MustPopEqualTo pops the next thing from the FIFO, waiting up to 1s for it to be available. If it is available within
// the timeout and equal to the expected value, it returns. Equality is tested with reflect.DeepEqual. If it is not
// available within the timeout, it triggers a fatal test failure. If it is not equal, it triggers a fatal test failure.
func (s *FIFO[I]) MustPopEqualTo(t Tester, expected I) {
	t.Helper()

	s.MustPopEqualToWithin(t, expected, 1*time.Second)
}

// MustPopAs pops the next thing from the FIFO, waiting up to 1s for it to be available. If it is available within the
// timeout, it attempts to set the target to be the value. If it is not available within the timeout, it triggers a
// fatal test failure. If it is not settable, it triggers a fatal test failure.
func (s *FIFO[I]) MustPopAs(t Tester, target I) {
	t.Helper()

	s.MustPopAsWithin(t, target, 1*time.Second)
}

// MustPopWithin pops the next thing from the FIFO, waiting up to the given duration for it to be available. If it is
// available, it returns the value. If it is not available within the timeout, it triggers a fatal test failure.
func (s *FIFO[I]) MustPopWithin(test Tester, d time.Duration) (next I) {
	test.Helper()

	var err error

	next, err = s.PopWithin(d)

	if err != nil {
		test.Fatal(err)
	}

	return
}

// MustPopEqualToWithin pops the next thing from the FIFO, waiting up to the given duration for it to be available. If
// it is available and equal to the expected value, it returns. Equality is tested with reflect.DeepEqual. If it is not
// available within the timeout, it triggers a fatal test failure. If it is not equal, it triggers a fatal test failure.
func (s *FIFO[I]) MustPopEqualToWithin(t Tester, expected I, d time.Duration) {
	t.Helper()

	err := s.PopEqualToWithin(expected, d)
	if err != nil {
		t.Fatal(err)
	}
}

// MustPopAsWithin pops the next thing from the FIFO, waiting up to the given duration for it to be available. If it is
// available, it attempts to set the target to be the value. If settable, it returns. If it is not available within the
// timeout, it triggers a fatal test failure. If it is not settable, it triggers a fatal test failure.
func (s *FIFO[I]) MustPopAsWithin(t Tester, target I, d time.Duration) {
	t.Helper()

	err := s.PopAsWithin(target, d)
	if err != nil {
		t.Fatal(err)
	}
}

// ConfirmClosed checks if the FIFO is closed and empty, waiting up to 1s. If it is not empty, it returns
// ErrChannelNotEmpty. If it is not closed within the timeout, it returns ErrChannelNotClosed.
func (s *FIFO[I]) ConfirmClosed() error {
	return s.ConfirmClosedWithin(1 * time.Second)
}

// ConfirmClosedWithin checks if the FIFO is closed and empty, waiting up to the given duration. If it is not empty, it
// returns ErrChannelNotEmpty. If it is not closed within the timeout, it returns ErrChannelNotClosed.
func (s *FIFO[I]) ConfirmClosedWithin(duration time.Duration) error {
	select {
	case value, ok := <-s.items:
		if ok {
			return fmt.Errorf("expected no more values in %s, but found %v: %w", s.name, value, ErrChannelNotEmpty)
		}
	case <-time.After(duration):
		return fmt.Errorf(
			"expected %s to be closed, but it was not after %v of waiting: %w",
			s.name,
			duration,
			ErrChannelNotClosed,
		)
	}

	return nil
}

// MustConfirmClosed checks if the FIFO is closed and empty, waiting up to 1s. If it is not empty, it triggers a fatal
// test failure. If it is not closed within the timeout, it triggers a fatal test failure.
func (s *FIFO[I]) MustConfirmClosed(t Tester) {
	t.Helper()

	s.MustConfirmClosedWithin(t, 1*time.Second)
}

// MustConfirmClosedWithin checks if the FIFO is closed and empty, waiting up to the given duration. If it is not empty,
// it triggers a fatal test failure. If it is not closed within the timeout, it triggers a fatal test failure.
func (s *FIFO[I]) MustConfirmClosedWithin(t Tester, d time.Duration) {
	t.Helper()

	err := s.ConfirmClosedWithin(d)
	if err != nil {
		t.Fatal(err)
	}
}

// Equal checks if the expected value (first arg) is equal to the actual value (second arg). Equality is tested with
// reflect.DeepEqual. If it is not equal, it returns ErrNotEqual.
func Equal[I any](expected, actual I) (err error) {
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected %T(%#v) but found %T(%#v) instead: %w", expected, expected, actual, actual, ErrNotEqual)
	}

	return nil
}

// MustEqual checks if the expected value (first arg) is equal to the actual value (second arg). Equality is tested with
// reflect.DeepEqual. If it is not equal, it triggers a fatal test failure.
func MustEqual[I any](t Tester, expected, actual I) {
	t.Helper()

	err := Equal(expected, actual)
	if err != nil {
		t.Fatal(err)
	}
}

// Helper Types for testing calls.
type CallWithNoArgsNoReturn struct{}

type CallWithNoArgs[R any] struct {
	ReturnOneShot *FIFO[R]
}

type CallWithNoReturn[A any] struct {
	Args A
}

type Call[A, R any] struct {
	Args          A
	ReturnOneShot *FIFO[R]
}

// ManageCall manages creating a call, setting its args, creating a return oneshot, recording the call, and returning
// the value popped off of the oneshot. The first type param here would be Call, but go won't allow that with ~.
func ManageCall[C ~struct {
	Args          A
	ReturnOneShot *FIFO[R]
}, A, R any](test Tester, calls *FIFO[any], args A,
) R {
	returnOneShot := NewOneShotFIFO[R]("return oneShot")

	// Turns out you can set fields on C, but not access them. So... setting it and pushing it here is fine, but setting
	// it first and then accessing the ReturnOneShot in the return statement is a compiler error. :shrug:
	calls.Push(C{
		Args:          args,
		ReturnOneShot: returnOneShot,
	})

	return returnOneShot.MustPop(test)
}

// ManageCallWithNoArgs manages creating a call, creating a return oneshot, recording the call, and returning the value
// popped off of the oneshot. The first type param here would be CallWithNoArgs, but go won't allow that with ~.
func ManageCallWithNoArgs[C ~struct {
	ReturnOneShot *FIFO[R]
}, R any](test Tester, calls *FIFO[any],
) R {
	returnOneShot := NewOneShotFIFO[R]("return oneShot")

	calls.Push(C{
		ReturnOneShot: returnOneShot,
	})

	return returnOneShot.MustPop(test)
}

// ManageCallWithNoArgsNoReturn manages creating a call and recording the call The first type param here would be
// CallWithNoArgsNoReturn, but go won't allow that with ~. This is not really helpful to reduce caller code, but it is
// helpful for consistency.
func ManageCallWithNoArgsNoReturn[C ~struct{}](calls *FIFO[any]) {
	calls.Push(C{})
}

type Tuple[V any] struct {
	Value V
	Err   error
}

func (t Tuple[V]) Unwrap() (V, error) {
	return t.Value, t.Err
}

// ManageCallWithNoReturn manages creating a call, setting its args, and recording the call. The first type param here
// would be CallWithNoReturn, but go won't allow that with ~.
func ManageCallWithNoReturn[C ~struct {
	Args A
}, A any](calls *FIFO[any], args A,
) {
	// Turns out you can set fields on C, but not access them. So... setting it and pushing it here is fine, but setting
	// it first and then accessing the ReturnOneShot in the return statement is a compiler error. :shrug:
	calls.Push(C{Args: args})
}

// ProxyCall proxies a call, creating oneshots for args and return, recording the call, pushing the args, and pulling &
// returning the value pulled for return. The first type param here would be Call, but go won't allow that with ~.
func ProxyCall(test Tester, calls *FIFO[AnyCall], name string, args ...any) []any {
	argsOneShot := NewOneShotFIFO[[]any]("args oneShot")
	returnOneShot := NewOneShotFIFO[[]any]("return oneShot")

	calls.Push(AnyCall{
		Name:    name,
		Args:    argsOneShot,
		Returns: returnOneShot,
	})

	argsOneShot.Push(args)

	return returnOneShot.MustPop(test)
}

type AnyCall struct {
	Name    string
	Args    *FIFO[[]any]
	Returns *FIFO[[]any]
}

func (ac *AnyCall) MustPullArgs(t Tester) []any {
	return ac.Args.MustPop(t)
}

func (ac *AnyCall) PushReturns(returns ...any) {
	ac.Returns.Push(returns)
}

// MustUnwrapTo1 checks that the input array has exactly 1 value of the specified generic type, and returns that.
func MustUnwrapTo1[R any](t Tester, returns []any) R {
	MustEqual(t, 1, len(returns))

	first, ok := returns[0].(R)
	if !ok {
		t.Fatal(fmt.Sprintf("failed type assertion of %#v to %T", returns[0], first))
	}

	return first
}

// MustUnwrapTo2 checks that the input array has exactly 2 values of the specified generic types, and returns them.
func MustUnwrapTo2[R1 any, R2 any](test Tester, returns []any) (first R1, second R2) {
	expectedNumValues := 2
	MustEqual(test, expectedNumValues, len(returns))

	var typeAssertionOk bool

	first, typeAssertionOk = returns[0].(R1)
	if !typeAssertionOk {
		test.Fatal(fmt.Sprintf("failed type assertion of item 0 (%#v) to %T", returns[0], first))
	}

	second, typeAssertionOk = returns[1].(R2)
	if !typeAssertionOk {
		test.Fatal(fmt.Sprintf("failed type assertion of item 1 (%#v) to %T", returns[1], second))
	}

	return first, second
}
