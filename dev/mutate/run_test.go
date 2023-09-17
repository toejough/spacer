//nolint:forcetypeassert // this is a test file, it's ok if type assertions panic
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

type simulator struct {
	deps     *runDeps
	callChan chan call
}

type call struct {
	name    string
	args    []any
	returns chan []any
}

func (sim *simulator) getCalled() call {
	select {
	case c := <-sim.callChan:
		return c
	case <-time.After(time.Second):
		panic("testing timeout waiting for a call")
	}
}

func (sim *simulator) shutdown() {
	close(sim.callChan)
}

var (
	errSimulatorNotShutDown     = errors.New("simulator was not shut down")
	errSimulatorShutdownTimeout = errors.New("simulator timed out waiting for shutdown")
)

func (sim *simulator) waitForShutdown() error {
	select {
	case thisCall, ok := <-sim.callChan:
		if !ok {
			// channel is closed
			return nil
		}

		return fmt.Errorf("had a call queued: %v: %w", thisCall, errSimulatorNotShutDown)
	case <-time.After(time.Second):
		return errSimulatorShutdownTimeout
	}
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

func newSimulator() *simulator {
	callChan := make(chan call)

	return &simulator{
		// TODO: separate the deps from the simulator.
		// TODO: hide the return channel stuff behind a function/method call.
		deps: &runDeps{
			printStarting: func(s string) {
				// no return channel for a func with no return
				callChan <- newCallNoReturn("printStarting", s)
			},
			printDoneWith: func(s string) {
				// no return channel for a func with no return
				callChan <- newCallNoReturn("printDoneWith", s)
			},
			pretest: func() bool {
				c := newCall("pretest")
				callChan <- c

				return (<-c.returns)[0].(bool)
			},
			testMutations: func() bool {
				c := newCall("testMutations")
				callChan <- c

				return (<-c.returns)[0].(bool)
			},
		},
		callChan: callChan,
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

func TestRunHappyPath(t *testing.T) {
	t.Parallel()

	// Given inputs
	sim := newSimulator()
	// and outputs
	var result bool

	// When the func is run
	go func() {
		result = run(sim.deps)
		sim.shutdown()
	}()

	// Then the start message is printed
	{
		actual := sim.getCalled()
		assertCalledNameIs(t, actual, "printStarting")
		assertArgsAre(t, actual, "Mutate")
	}

	// Then the pretest is run
	{
		actual := sim.getCalled()
		assertCalledNameIs(t, actual, "pretest")
		actual.injectReturn(true)
	}

	// Then the mutation testing is run
	{
		actual := sim.getCalled()
		assertCalledNameIs(t, actual, "testMutations")
		actual.injectReturn(true)
	}

	// Then the done message is printed
	{
		actual := sim.getCalled()
		assertCalledNameIs(t, actual, "printDoneWith")
		assertArgsAre(t, actual, "Mutate")
	}

	// Then expect that the simulator is done
	if err := sim.waitForShutdown(); err != nil {
		t.Fatalf("the simulator is not done yet at the end of the test: %s", err)
	}

	// Then the result is true
	if result != true {
		t.Fatal("The result was false")
	}
}
