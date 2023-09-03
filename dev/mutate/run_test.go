//nolint:forcetypeassert // this is a test file, it's ok if type assertions panic
package main

// The main idea for all the unit tests is to test the behavior we care about _at this level_.
// This means we validate the calls to dependencies _at this level_ (critically, _not_ subdependency calls).
// Leave "and now xyz is happening" testing to the thing that is making it happen.
// For example, for "run", we do _not_ care where the test command is coming from, how it is run, or how its output is
// conveyed.

import (
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

const printStartingName = "printStarting"

func newSimulator() *simulator {
	callChan := make(chan call)

	return &simulator{
		deps: &runDeps{
			printStarting: func(s string) {
				returnChan := make(chan []any)
				callChan <- call{name: printStartingName, args: []any{s}, returns: returnChan}
				<-returnChan
			},
			printDoneWith: func(string) {},
			pretest:       func() bool { return false },
			testMutations: func() bool { return false },
		},
		callChan: callChan,
	}
}

func TestRunHappyPath(t *testing.T) {
	t.Parallel()

	// Given inputs
	sim := newSimulator()
	// and outputs

	// When the func is run
	go func() {
		run(sim.deps)
	}()

	// Then the start message is printed
	{
		expectedName := printStartingName
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
	// Then the mutation testing is run
	// Then the done message is printed
} //nolint:wsl // these are definitely todo-style comments & I want them here for now
