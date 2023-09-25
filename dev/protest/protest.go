// Package protest provides procedure testing functionality.
package protest

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

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

type CallRelay struct {
	callChan chan Call
}

type Call struct {
	name    string
	args    []any
	returns chan []any
}

type RelayReader interface {
	Get() callTester
	WaitForShutdown(time.Duration) error
}

type RelayWriter interface {
	Put(Call) returnReader
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

func NewCall(name string, args ...any) Call {
	return Call{name: name, args: args, returns: make(chan []any)}
}

func NewCallNoReturn(name string, args ...any) Call {
	return Call{name: name, args: args, returns: nil}
}

func (c Call) Name() string {
	return c.name
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
		panic("timed out waiting for " + c.name + " to read the injected return values")
	}
}

func (c Call) FillReturns(returnPointers ...any) {
	returnValues := <-c.returns
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

func (rt *RelayTester) AssertNextCallIs(message string, args ...any) callTester {
	rt.T.Helper()
	return AssertNextCallIs(rt.T, rt.Relay, message, args...)
}

func (rt *RelayTester) AssertRelayShutsDownWithin(d time.Duration) {
	rt.T.Helper()
	AssertRelayShutsDownWithin(rt.T, rt.Relay, d)
}
