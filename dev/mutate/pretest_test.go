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

type pretestDeps interface {
	printStarting(string) func(string)
	fetchPretestCommand() []string
	runSubprocess([]string)
}

type mockPretestDeps struct{ relay *callRelay }

func (d *mockPretestDeps) printStarting(message string) func(string) {
	var returnFunc func(string)
	// TODO: relay.putNewCall() ?? no... this is not something we're doing to relay, it's something we're doing with it...
	// make it putNewCall(relay, ...)
	// TODO: make directional relay, so that you can't both put & get calls from the same object/interface
	d.relay.putCall(newCall("printStarting", message)).fillReturns(&returnFunc)

	return returnFunc
}

func (d *mockPretestDeps) fetchPretestCommand() []string {
	var c []string

	d.relay.putCall(newCall("fetchPretestCommand")).fillReturns(&c)

	return c
}

func (d *mockPretestDeps) runSubprocess(command []string) {
	d.relay.putCall(newCall("runSubprocess", command))
}

// TODO: move non-test functions out to their own file.
func pretest(deps pretestDeps) bool {
	done := deps.printStarting("Pretest")
	defer done("Success")

	command := deps.fetchPretestCommand()
	deps.runSubprocess(command)

	return true
}

func newPretestDeps(relay *callRelay) *mockPretestDeps {
	return &mockPretestDeps{relay: relay}
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
	// TODO: create assertNextCallIs(t, relay, ...)
	// TODO: create tester. newTester(t, relay).assertNextCallIs(...)
	assertCallIs(t, relay.getCall(), "printStarting", "Pretest").injectReturn(mockDoneFunc)
	// Then the pretest is fetched
	assertCallIs(t, relay.getCall(), "fetchPretestCommand").injectReturn(pretestCommand)
	// Then the pretest command is run
	assertCallIs(t, relay.getCall(), "runSubprocess", pretestCommand)
	// Then the done message is printed
	assertCallIs(t, relay.getCall(), "printDone", "Success")

	// Then the relay is shut down
	assertRelayShutsDownWithin(t, relay, time.Second)

	// Then the functin passed
	if !passed {
		t.Fatal("the pretest function failed unexpectedly")
	}
}
