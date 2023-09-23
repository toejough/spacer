package main

// The main idea for all the unit tests is to test the behavior we care about
// _at this level_. This means we validate the calls to dependencies _at this
// level_ (critically, _not_ subdependency calls). Leave "and now xyz is
// happening" testing to the thing that is making it happen. For example, for
// "run", we do _not_ care where the pretest command is coming from, how it is
// run, or how its output is conveyed.

import (
	"testing"
	"time"
)

func newPretestDeps(relay *callRelay) *runDeps {
	return &pretestDeps{}
}

func TestPretestHappyPath(t *testing.T) {
	t.Parallel()

	// Given inputs
	relay := newCallRelay()
	deps := newPretestDeps(relay)
	mockDoneFunc := func(message string) { relay.putCall(newCallNoReturn("printDone", message)) }
	pretestCommand := []string{"this", "is", "a", "test", "command"}
	// and outputs
	passed := false

	// When the func is run
	go func() {
		passed = pretest(deps)

		relay.shutdown()
	}()

	// Then the start message is printed
	assertCallIs(t, relay.getCall(), "printStarting", "Pretest").injectReturn(mockDoneFunc)
	// Then the pretest is fetched
	assertCallIs(t, relay.getCall(), "fetchPretestCommand").injectReturn(pretestCommand)
	// Then the pretest command is run
	assertCallIs(t, relay.getCall(), "runSubprocess", pretestCommand).injectReturn(true)
	// Then the done message is printed
	assertCallIs(t, relay.getCall(), "printDone", "Success")
	// Then the program exits with 0
	assertCallIs(t, relay.getCall(), "exit", 0)

	// Then the relay is shut down
	assertRelayShutsDownWithin(t, relay, time.Second)

	// Then the functin passed
	if !passed {
		t.Fatal("the pretest function failed unexpectedly")
	}
}
