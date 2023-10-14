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

func (rd *mockRunDeps) pretest() bool {
	var b bool

	rd.relay.PutCall(rd.pretest).FillReturns(&b)

	return b
}

func (rd *mockRunDeps) testMutations() bool {
	var success bool

	rd.relay.PutCall(rd.testMutations).FillReturns(&success)

	return success
}

func (rd *mockRunDeps) exit(code int) {
	rd.relay.PutCallNoReturn(rd.exit, code)
}

func TestRunHappyPath(t *testing.T) {
	t.Parallel()

	// Given test needs
	relay := protest.NewCallRelay()
	tester := &protest.RelayTester{T: t, Relay: relay}
	// Given inputs
	deps := &mockRunDeps{relay: relay}

	// When the func is run
	tester.Start(run, deps)

	// Then the pretest is run
	tester.AssertNextCallIs(deps.pretest).InjectReturns(true)
	// Then the mutation testing is run
	tester.AssertNextCallIs(deps.testMutations).InjectReturns(true)
	// Then the program exits with 0
	tester.AssertNextCallIs(deps.exit, 0)

	// Then the relay is shut down
	tester.AssertDoneWithin(time.Second)
}

func TestRunPretestFailure(t *testing.T) {
	t.Parallel()

	// Given test needs
	relay := protest.NewCallRelay()
	tester := &protest.RelayTester{T: t, Relay: relay}
	// Given inputs
	deps := &mockRunDeps{relay: relay}

	// When the func is run
	tester.Start(run, deps)

	// Then the pretest is run
	tester.AssertNextCallIs(deps.pretest).InjectReturns(false)
	// Then the program exits with 1
	tester.AssertNextCallIs(deps.exit, 1)

	// Then the relay is shut down
	tester.AssertDoneWithin(time.Second)
}

func TestRunMutationFailure(t *testing.T) {
	t.Parallel()

	// Given test needs
	relay := protest.NewCallRelay()
	tester := &protest.RelayTester{T: t, Relay: relay}
	// Given inputs
	deps := &mockRunDeps{relay: relay}

	// When the func is run
	tester.Start(run, deps)

	// Then the pretest is run
	tester.AssertNextCallIs(deps.pretest).InjectReturns(true)
	// Then the mutation testing is run
	tester.AssertNextCallIs(deps.testMutations).InjectReturns(false)
	// Then the program exits with 1
	tester.AssertNextCallIs(deps.exit, 1)

	// Then the relay is shut down
	tester.AssertDoneWithin(time.Second)
}
