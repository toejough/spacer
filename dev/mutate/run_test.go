//nolint:forcetypeassert // this is a test file, it's ok if type assertions panic
package main

// The main idea for all the unit tests is to test the behavior we care about _at this level_.
// This means we validate the calls to dependencies _at this level_ (critically, _not_ subdependency calls).
// Leave "and now xyz is happening" testing to the thing that is making it happen.
// For example, for "run", we do _not_ care where the test command is coming from, how it is run, or how its output is
// conveyed.

import (
	"errors"
	"fmt"
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

func newSimulator() *simulator {
	callChan := make(chan call)

	return &simulator{
		deps: &runDeps{
			printStarting: func(s string) {
				returnChan := make(chan []any)
				callChan <- call{name: "printStarting", args: []any{s}, returns: returnChan}
				<-returnChan
			},
			printDoneWith: func(s string) {
				returnChan := make(chan []any)
				callChan <- call{name: "printDoneWith", args: []any{s}, returns: returnChan}
				<-returnChan
			},
			pretest: func() bool {
				returnChan := make(chan []any)
				callChan <- call{name: "pretest", args: []any{}, returns: returnChan}

				return (<-returnChan)[0].(bool)
			},
			testMutations: func() bool {
				returnChan := make(chan []any)
				callChan <- call{name: "testMutations", args: []any{}, returns: returnChan}

				return (<-returnChan)[0].(bool)
			},
		},
		callChan: callChan,
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
		expectedName := "printStarting"
		actual := sim.getCalled()
		if actual.name != expectedName {
			t.Fatalf("the called function was expected to be %s, but was %s instead", expectedName, actual.name)
		}
		expectedArgs := "Mutate"
		if actual.args[0].(string) != expectedArgs {
			t.Fatalf("the function %s was expected to be called with %s but was called with %s",
				actual.name, expectedArgs, actual.args,
			)
		}
		actual.returns <- nil
	}

	// Then the pretest is run
	{
		expectedName := "pretest"
		actual := sim.getCalled()
		if actual.name != expectedName {
			t.Fatalf("the called function was expected to be %s, but was %s instead", expectedName, actual.name)
		}
		expectedArgs := []any{}
		if len(actual.args) != len(expectedArgs) {
			t.Fatalf("the function %s was expected to be called with %d args but was called with %d args instead",
				actual.name, len(expectedArgs), len(actual.args),
			)
		}
		actual.returns <- []any{true}
	}

	// Then the mutation testing is run
	{
		expectedName := "testMutations"
		actual := sim.getCalled()
		if actual.name != expectedName {
			t.Fatalf("the called function was expected to be %s, but was %s instead", expectedName, actual.name)
		}
		expectedArgs := []any{}
		if len(actual.args) != len(expectedArgs) {
			t.Fatalf("the function %s was expected to be called with %d args but was called with %d args instead",
				actual.name, len(expectedArgs), len(actual.args),
			)
		}
		actual.returns <- []any{true}
	}

	// Then the done message is printed
	{
		expectedName := "printDoneWith"
		actual := sim.getCalled()
		if actual.name != expectedName {
			t.Fatalf("the called function was expected to be %s, but was %s instead", expectedName, actual.name)
		}
		expectedArgs := "Mutate"
		if actual.args[0].(string) != expectedArgs {
			t.Fatalf("the function %s was expected to be called with %s but was called with %s",
				actual.name, expectedArgs, actual.args,
			)
		}
		actual.returns <- nil
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
