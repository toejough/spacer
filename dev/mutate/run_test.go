package main

// The main idea for all the unit tests is to test the behavior we care about
// _at this level_. This means we validate the calls to dependencies _at this
// level_ (critically, _not_ subdependency calls). Leave "and now xyz is
// happening" testing to the thing that is making it happen. For example, for
// "run", we do _not_ care where the pretest command is coming from, how it is
// run, or how its output is conveyed.

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type callRelay struct {
	callChan chan call
}

type call struct {
	name    string
	args    []any
	returns chan []any
}

func (cr *callRelay) getCall() call {
	select {
	case c := <-cr.callChan:
		return c
	case <-time.After(time.Second):
		panic("testing timeout waiting for a call")
	}
}

func (cr *callRelay) shutdown() {
	close(cr.callChan)
}

var (
	errCallRelayNotShutDown     = errors.New("call relay was not shut down")
	errCallRelayShutdownTimeout = errors.New("call relay timed out waiting for shutdown")
)

func (cr *callRelay) waitForShutdown(waitTime time.Duration) error {
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

func (cr *callRelay) putCall(c call) call {
	cr.callChan <- c
	return c
}

func newCall(name string, args ...any) call {
	return call{name: name, args: args, returns: make(chan []any)}
}

func newCallNoReturn(name string, args ...any) call {
	return call{name: name, args: args, returns: nil}
}

func (c call) injectReturn(returnValues ...any) {
	c.returns <- returnValues
}

func (c call) fillReturns(returnPointers ...any) {
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

func newCallRelay() *callRelay {
	return &callRelay{callChan: make(chan call)}
}

func newDeps(relay *callRelay) *runDeps {
	return &runDeps{
		printStarting: func(s string) {
			relay.putCall(newCallNoReturn("printStarting", s))
		},
		printDoneWith: func(s string) {
			relay.putCall(newCallNoReturn("printDoneWith", s))
		},
		pretest: func() bool {
			var b bool
			relay.putCall(newCall("pretest")).fillReturns(&b)

			return b
		},
		testMutations: func() bool {
			var b bool
			relay.putCall(newCall("testMutations")).fillReturns(&b)

			return b
		},
	}
}

func assertCalledNameIs(t *testing.T, c call, expectedName string) {
	t.Helper()

	if c.name != expectedName {
		t.Fatalf("the called function was expected to be %s, but was %s instead", expectedName, c.name)
	}
}

func assertArgsAre(t *testing.T, c call, expectedArgs ...any) {
	t.Helper()

	if !reflect.DeepEqual(c.args, expectedArgs) {
		t.Fatalf("the function %s was expected to be called with %#v but was called with %#v",
			c.name, expectedArgs, c.args,
		)
	}
}

func assertCallIs(t *testing.T, c call, name string, expectedArgs ...any) call {
	t.Helper()
	assertCalledNameIs(t, c, name)
	assertArgsAre(t, c, expectedArgs...)

	return c
}

func assertRelayShutsDownWithin(t *testing.T, relay *callRelay, waitTime time.Duration) {
	t.Helper()

	if err := relay.waitForShutdown(waitTime); err != nil {
		t.Fatalf("the simulator is not done yet: %s", err)
	}
}

func TestRunHappyPath(t *testing.T) {
	t.Parallel()

	// Given inputs
	relay := newCallRelay()
	deps := newDeps(relay)
	// and outputs
	var result bool

	// When the func is run
	go func() {
		result = run(deps)

		relay.shutdown()
	}()

	// Then the start message is printed
	assertCallIs(t, relay.getCall(), "printStarting", "Mutate")
	// Then the pretest is run
	assertCallIs(t, relay.getCall(), "pretest").injectReturn(true)
	// Then the mutation testing is run
	assertCallIs(t, relay.getCall(), "testMutations").injectReturn(true)
	// Then the done message is printed
	assertCallIs(t, relay.getCall(), "printDoneWith", "Mutate")

	// Then the relay is shut down
	assertRelayShutsDownWithin(t, relay, time.Second)

	// Then the result is true
	if result != true {
		t.Fatal("The result was false")
	}
}

func TestRunPretestFailure(t *testing.T) {
	t.Parallel()

	// Given inputs
	relay := newCallRelay()
	deps := newDeps(relay)
	// and outputs
	var result bool

	// When the func is run
	go func() {
		result = run(deps)

		relay.shutdown()
	}()

	// Then the start message is printed
	assertCallIs(t, relay.getCall(), "printStarting", "Mutate")
	// Then the pretest is run
	assertCallIs(t, relay.getCall(), "pretest").injectReturn(false)
	// Then the done message is printed
	assertCallIs(t, relay.getCall(), "printDoneWith", "Mutate")

	// Then the relay is shut down
	assertRelayShutsDownWithin(t, relay, time.Second)

	// Then the result is true
	if result != false {
		t.Fatal("The result was unexpectedly true")
	}
}

func TestRunMutationFailure(t *testing.T) {
	t.Parallel()

	// Given inputs
	relay := newCallRelay()
	deps := newDeps(relay)
	// and outputs
	var result bool

	// When the func is run
	go func() {
		result = run(deps)

		relay.shutdown()
	}()

	// Then the start message is printed
	assertCallIs(t, relay.getCall(), "printStarting", "Mutate")
	// Then the pretest is run
	assertCallIs(t, relay.getCall(), "pretest").injectReturn(true)
	// Then the mutation testing is run
	assertCallIs(t, relay.getCall(), "testMutations").injectReturn(false)
	// Then the done message is printed
	assertCallIs(t, relay.getCall(), "printDoneWith", "Mutate")

	// Then the relay is shut down
	assertRelayShutsDownWithin(t, relay, time.Second)

	// Then the result is true
	if result != false {
		t.Fatal("The result was unexpectedly true")
	}
}
