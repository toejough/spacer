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

	"pgregory.net/rapid"
)

type mockPretestDeps struct{ relay *protest.CallRelay }

func (d *mockPretestDeps) printStarting(message string) func(string) {
	var returnFunc func(string)

	d.relay.PutCall(d.printStarting, message).FillReturns(&returnFunc)

	return returnFunc
}

func (d *mockPretestDeps) fetchPretestCommand() []string {
	var c []string

	d.relay.PutCall(d.fetchPretestCommand).FillReturns(&c)

	return c
}

func (d *mockPretestDeps) runSubprocess(command []string) bool {
	var b bool

	d.relay.PutCall(d.runSubprocess, command).FillReturns(&b)

	return b
}

func (d *mockPretestDeps) printDone(message string) {
	d.relay.PutCallNoReturn(d.printDone, message)
}

func newPretestDeps(relay *protest.CallRelay) *mockPretestDeps {
	return &mockPretestDeps{relay: relay}
}

func rapidPretestHappyPath(rapidTester *rapid.T) {
	// Given test needs
	relay := protest.NewCallRelay()
	tester := &protest.RelayTester{T: rapidTester, Relay: relay}
	// Given inputs
	deps := newPretestDeps(relay)
	pretestCommand := rapid.SliceOf(rapid.String()).Draw(rapidTester, "pretestCommand")
	// and outputs
	passed := false

	// When the func is run
	go func() {
		passed = pretest(deps)

		relay.Shutdown()
	}()

	// Then the start message is printed
	tester.AssertNextCallIs(deps.printStarting, "Pretest").InjectReturns(deps.printDone)
	// Then the pretest is fetched
	tester.AssertNextCallIs(deps.fetchPretestCommand).InjectReturns(pretestCommand)
	// Then the pretest command is run
	tester.AssertNextCallIs(deps.runSubprocess, pretestCommand).InjectReturns(true)
	// Then the done message is printed
	tester.AssertNextCallIs(deps.printDone, "Success")

	// Then the relay is shut down
	tester.AssertRelayShutsDownWithin(time.Second)

	// Then the functin passed
	if !passed {
		rapidTester.Fatal("the pretest function failed unexpectedly")
	}
}

func FuzzPretestHappyPath(f *testing.F) {
	f.Fuzz(rapid.MakeFuzz(rapidPretestHappyPath))
}

func TestPretestHappyPath(t *testing.T) {
	t.Parallel()
	rapid.Check(t, rapidPretestHappyPath)
}

func TestPretestSubprocessFail(t *testing.T) {
	t.Parallel()
	rapid.Check(t, func(rapidTester *rapid.T) {
		// Given test needs
		relay := protest.NewCallRelay()
		tester := &protest.RelayTester{T: t, Relay: relay}
		// Given inputs
		deps := newPretestDeps(relay)
		pretestCommand := rapid.SliceOf(rapid.String()).Draw(rapidTester, "pretestCommand")
		// and outputs
		passed := false

		// When the func is run
		go func() {
			passed = pretest(deps)

			relay.Shutdown()
		}()

		// Then the start message is printed
		tester.AssertNextCallIs(deps.printStarting, "Pretest").InjectReturns(deps.printDone)
		// Then the pretest is fetched
		tester.AssertNextCallIs(deps.fetchPretestCommand).InjectReturns(pretestCommand)
		// Then the pretest command is run
		tester.AssertNextCallIs(deps.runSubprocess, pretestCommand).InjectReturns(false)
		// Then the done message is printed
		tester.AssertNextCallIs(deps.printDone, "Failure")

		// Then the relay is shut down
		tester.AssertRelayShutsDownWithin(time.Second)

		// Then the function failed
		if passed {
			t.Fatal("the pretest function passed unexpectedly")
		}
	})
}
