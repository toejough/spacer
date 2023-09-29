package main

// The main idea for all the unit tests is to test the behavior we care about
// _at this level_. This means we validate the calls to dependencies _at this
// level_ (critically, _not_ subdependency calls). Leave "and now xyz is
// happening" testing to the thing that is making it happen. For example, for
// "run", we do _not_ care where the pretest command is coming from, how it is
// run, or how its output is conveyed.

import (
	"spacer/dev/protest"
	"testing"
	"time"
)

type mockPretestDeps struct{ relay protest.RelayWriter }

// TODO: a relay function that creates a new call with the name already read from the calling method.
func (d *mockPretestDeps) printStarting(message string) func(string) {
	var returnFunc func(string)

	d.relay.Put(protest.NewCall("printStarting", message)).FillReturns(&returnFunc)

	return returnFunc
}

func (d *mockPretestDeps) fetchPretestCommand() []string {
	var c []string

	d.relay.Put(protest.NewCall("fetchPretestCommand")).FillReturns(&c)

	return c
}

func (d *mockPretestDeps) runSubprocess(command []string) {
	d.relay.Put(protest.NewCall("runSubprocess", command))
}

func newPretestDeps(relay protest.RelayWriter) *mockPretestDeps {
	return &mockPretestDeps{relay: relay}
}

func TestPretestHappyPath(t *testing.T) {
	t.Parallel()

	// Given test needs
	relay := protest.NewCallRelay()
	tester := &protest.RelayTester{T: t, Relay: relay}
	// Given inputs
	deps := newPretestDeps(relay)
	mockDoneFunc := func(message string) { relay.Put(protest.NewCallNoReturn("printDone", message)) }
	pretestCommand := []string{"this", "is", "a", "test", "command"}
	// and outputs
	passed := false

	// When the func is run
	go func() {
		passed = pretest(deps)

		relay.Shutdown()
	}()

	// Then the start message is printed
	// TODO: a protest function to read the name from a given function
	// TODO: update tester function to use the function instead of a name
	tester.AssertNextCallIs("printStarting", "Pretest").InjectReturn(mockDoneFunc)
	// Then the pretest is fetched
	// TODO: do property testing for the command returned
	tester.AssertNextCallIs("fetchPretestCommand").InjectReturn(pretestCommand)
	// Then the pretest command is run
	// TODO: add a test for when the subprocess command fails
	tester.AssertNextCallIs("runSubprocess", pretestCommand)
	// Then the done message is printed
	tester.AssertNextCallIs("printDone", "Success")

	// Then the relay is shut down
	tester.AssertRelayShutsDownWithin(time.Second)

	// Then the functin passed
	if !passed {
		t.Fatal("the pretest function failed unexpectedly")
	}
}
