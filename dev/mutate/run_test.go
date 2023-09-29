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

type mockRunDeps struct {
	relay *protest.CallRelay
}

func (rd *mockRunDeps) printStarting(s string) func(string) {
	var f func(string)

	rd.relay.Put(protest.NewCall("printStarting", s)).FillReturns(&f)

	return f
}

func (rd *mockRunDeps) pretest() bool {
	var b bool

	rd.relay.Put(protest.NewCall("pretest")).FillReturns(&b)

	return b
}

func (rd *mockRunDeps) testMutations() bool {
	var success bool

	// TODO: is there a way to grab the name like testify.mock does?
	// TODO keep this filling, or do type assertions on returns(0) the way testify.mock does?
	rd.relay.Put(protest.NewCall("testMutations")).FillReturns(&success)

	return success
}

func (rd *mockRunDeps) exit(code int) {
	rd.relay.Put(protest.NewCallNoReturn("exit", code))
}

// TODO: implement the same tests with https://github.com/stretchr/testify/issues/741 and see how that feels.
func TestRunHappyPath(t *testing.T) {
	t.Parallel()

	// Given inputs
	relay := protest.NewCallRelay()
	deps := &mockRunDeps{relay: relay}
	mockDoneFunc := func(message string) { relay.Put(protest.NewCallNoReturn("printDone", message)) }

	// When the func is run
	go func() {
		run(deps)

		relay.Shutdown()
	}()

	// Then the start message is printed
	protest.AssertNextCallIs(t, relay, "printStarting", "Mutate").InjectReturn(mockDoneFunc)
	// Then the pretest is run
	protest.AssertNextCallIs(t, relay, "pretest").InjectReturn(true)
	// Then the mutation testing is run
	protest.AssertNextCallIs(t, relay, "testMutations").InjectReturn(true)
	// Then the done message is printed
	protest.AssertNextCallIs(t, relay, "printDone", "Success")
	// Then the program exits with 0
	protest.AssertNextCallIs(t, relay, "exit", 0)

	// Then the relay is shut down
	protest.AssertRelayShutsDownWithin(t, relay, time.Second)
}

func TestRunPretestFailure(t *testing.T) {
	t.Parallel()

	// Given inputs
	relay := protest.NewCallRelay()
	deps := &mockRunDeps{relay: relay}
	mockDoneFunc := func(message string) { relay.Put(protest.NewCallNoReturn("printDone", message)) }

	// When the func is run
	go func() {
		run(deps)

		relay.Shutdown()
	}()

	// Then the start message is printed
	protest.AssertNextCallIs(t, relay, "printStarting", "Mutate").InjectReturn(mockDoneFunc)
	// Then the pretEst is run
	protest.AssertNextCallIs(t, relay, "pretest").InjectReturn(false)
	// Then the done message is printed
	protest.AssertNextCallIs(t, relay, "printDone", "Failure")
	// Then the program exits with 1
	protest.AssertNextCallIs(t, relay, "exit", 1)

	// Then the relay is shut down
	protest.AssertRelayShutsDownWithin(t, relay, time.Second)
}

func TestRunMutationFailure(t *testing.T) {
	t.Parallel()

	// Given inputs
	relay := protest.NewCallRelay()
	deps := &mockRunDeps{relay: relay}
	mockDoneFunc := func(message string) { relay.Put(protest.NewCallNoReturn("printDone", message)) }

	// When the func is run
	go func() {
		run(deps)

		relay.Shutdown()
	}()

	// Then the start message is printed
	protest.AssertNextCallIs(t, relay, "printStarting", "Mutate").InjectReturn(mockDoneFunc)
	// Then the pretest is run
	protest.AssertNextCallIs(t, relay, "pretest").InjectReturn(true)
	// Then the mutation testing is run
	protest.AssertNextCallIs(t, relay, "testMutations").InjectReturn(false)
	// Then the done message is printed
	protest.AssertNextCallIs(t, relay, "printDone", "Failure")
	// Then the program exits with 1
	protest.AssertNextCallIs(t, relay, "exit", 1)

	// Then the relay is shut down
	protest.AssertRelayShutsDownWithin(t, relay, time.Second)
}
