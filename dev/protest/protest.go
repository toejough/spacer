// Package protest provides procedure testing functionality.
// Why not https://github.com/stretchr/testify/blob/master/README.md#mock-package?
// You only get to specify simple call/return behavior, with no guarantees about ordering, and you need to unset
// handlers for repeated calls for the same function.
// On the other hand, there's https://github.com/stretchr/testify/issues/741.  Is this necessary?
// maybe this whole suite is pointless?
package protest

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
)

type (
	CallRelay struct {
		callChan chan Call
	}
	Call struct {
		function Function
		args     []any
		returns  chan []any
	}
	RelayReader interface {
		Get() callTester
		WaitForShutdown(time.Duration) error
	}
	RelayWriter interface {
		Put(Call) returnReader
		PutCall(Function, ...any) returnReader
		PutCallNoReturn(Function, ...any) returnReader
	}
)

// Public helpers.
func AssertNextCallIs(t *testing.T, r RelayReader, name string, expectedArgs ...any) callTester {
	t.Helper()

	c := r.Get()
	assertCalledNameIs(t, c, name)
	assertArgsAre(t, c, expectedArgs...)

	return c
}

func AssertRelayShutsDownWithin(t *testing.T, relay RelayReader, waitTime time.Duration) {
	t.Helper()

	if err := relay.WaitForShutdown(waitTime); err != nil {
		t.Fatalf("the relay has not shut down yet: %s", err)
	}
}

// Private helpers.
func assertCalledNameIs(t *testing.T, c callReader, expectedName string) {
	t.Helper()

	if c.Name() != expectedName {
		t.Fatalf("the called function was expected to be %s, but was %s instead", expectedName, c.Name())
	}
}

func assertArgsAre(t *testing.T, theCall callReader, expectedArgs ...any) {
	t.Helper()

	if theCall.Args() == nil && expectedArgs != nil {
		t.Fatalf(
			"the function %s was expected to be called with %#v, but was called without args",
			theCall.Name(),
			expectedArgs,
		)
	}

	if theCall.Args() != nil && expectedArgs == nil {
		t.Fatalf(
			"the function %s was expected to be called without args, but was called with %#v",
			theCall.Name(),
			theCall.Args(),
		)
	}

	if !reflect.DeepEqual(theCall.Args(), expectedArgs) {
		t.Fatalf("the function %s was expected to be called with %#v but was called with %#v",
			theCall.Name(), expectedArgs, theCall.Args(),
		)
	}
}

func (cr *CallRelay) Get() callTester {
	select {
	case c, ok := <-cr.callChan:
		if !ok {
			panic("expected a call, but the relay was already shut down")
		}

		return c
	case <-time.After(time.Second):
		panic("testing timeout waiting for a call")
	}
}

func (cr *CallRelay) Shutdown() {
	close(cr.callChan)
}

var (
	errCallRelayNotShutDown     = errors.New("call relay was not shut down")
	errCallRelayShutdownTimeout = errors.New("call relay timed out waiting for shutdown")
)

func (cr *CallRelay) WaitForShutdown(waitTime time.Duration) error {
	select {
	case thisCall, ok := <-cr.callChan:
		if !ok {
			// channel is closed
			return nil
		}

		return fmt.Errorf("had a call queued: %v: %w", thisCall, errCallRelayNotShutDown)
	case <-time.After(waitTime):
		return errCallRelayShutdownTimeout
	}
}

type callReader interface {
	Name() string
	Args() []any
}

type returnWriter interface {
	InjectReturn(...any)
}

type callTester interface {
	callReader
	returnWriter
}

type returnReader interface {
	FillReturns(...any)
}

func (cr *CallRelay) Put(c Call) returnReader {
	cr.callChan <- c
	return c
}

type Function any

func getFuncName(f Function) string {
	// docs say to use UnsafePointer explicitly instead of Pointer()
	// https://pkg.go.dev/reflect@go1.21.1#Value.Pointer
	return runtime.FuncForPC(uintptr(reflect.ValueOf(f).UnsafePointer())).Name()
}

func (cr *CallRelay) PutCall(f Function, args ...any) returnReader {
	return cr.Put(NewCall(f, args...))
}

// TODO: try to return concrete types.
func (cr *CallRelay) PutCallNoReturn(f Function, args ...any) returnReader {
	return cr.Put(NewCallNoReturn(f, args...))
}

func panicIfNotFunc(evaluate Function, from Function) {
	kind := reflect.ValueOf(evaluate).Kind()
	if kind != reflect.Func {
		panic(fmt.Sprintf("must pass a function as the first argument to %s. received a %s instead.",
			getFuncName(from),
			kind.String(),
		))
	}
}

func NewCall(f Function, args ...any) Call {
	panicIfNotFunc(f, NewCall)
	return Call{function: f, args: args, returns: make(chan []any)}
}

func NewCallNoReturn(f Function, args ...any) Call {
	panicIfNotFunc(f, NewCallNoReturn)
	return Call{function: f, args: args, returns: nil}
}

func (c Call) Name() string {
	return getFuncName(c.function)
}

func (c Call) Args() []any {
	return c.args
}

func (c Call) InjectReturn(returnValues ...any) {
	if c.returns == nil {
		panic("cannot inject a return on a call with no returns")
	}
	select {
	case c.returns <- returnValues:
		return
	case <-time.After(1 * time.Second):
		panic("timed out waiting for " + c.Name() + " to read the injected return values")
	}
}

func (c Call) FillReturns(returnPointers ...any) {
	returnValues := <-c.returns
	// TODO: callout in the docs: it's ok to panic if the test is written wrong. Not ok to panic if the test is just failing.
	// TODO: document as a win over testify.Mocked: we fail if the number of returns doesn't match
	// TODO: Can we figure out the number of returns from the caller's function signature?
	for index := range returnValues {
		// USEFUL SNIPPETS FROM JSON.UNMARSHAL
		// if rv.Kind() != reflect.Pointer || rv.IsNil() {
		// 	return &InvalidUnmarshalError{reflect.TypeOf(v)}
		// }
		// v.Set(reflect.ValueOf(oi))
		rv := reflect.ValueOf(returnPointers[index])
		if rv.Kind() != reflect.Pointer || rv.IsNil() {
			panic("cannot fill value into non-pointer")
		}
		// Use Elem instead of directly using Set for setting pointers
		rv.Elem().Set(reflect.ValueOf(returnValues[index]))
	}
}

func NewCallRelay() *CallRelay {
	return &CallRelay{callChan: make(chan Call)}
}

type RelayTester struct {
	T     *testing.T
	Relay RelayReader
}

// TODO: can we know the number of args & check that here?
func (rt *RelayTester) AssertNextCallIs(f Function, args ...any) callTester {
	rt.T.Helper()
	panicIfNotFunc(f, AssertNextCallIs)

	return AssertNextCallIs(rt.T, rt.Relay, getFuncName(f), args...)
}

func (rt *RelayTester) AssertRelayShutsDownWithin(d time.Duration) {
	rt.T.Helper()
	AssertRelayShutsDownWithin(rt.T, rt.Relay, d)
}
