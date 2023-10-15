package main

// The main idea for all the unit tests is to test the behavior we care about
// _at this level_. This means we validate the calls to dependencies _at this
// level_ (critically, _not_ subdependency calls). Leave "and now xyz is
// happening" testing to the thing that is making it happen. For example, for
// "run", we do _not_ care where the pretest command is coming from, how it is
// run, or how its output is conveyed.

import (
	"errors"
	"spacer/dev/protest"
	"testing"
	"time"

	"pgregory.net/rapid"
)

type mockPretestDeps struct{ relay *protest.CallRelay }

func (d *mockPretestDeps) fetchPretestCommand() []string {
	var c []string

	d.relay.PutCall(d.fetchPretestCommand).FillReturns(&c)

	return c
}

func (d *mockPretestDeps) runSubprocess(command []string) error {
	var e error

	d.relay.PutCall(d.runSubprocess, command).FillReturns(&e)

	return e
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
	// Given outputs
	result := true

	// When the func is run
	tester.Start(pretest, deps)

	// Then the pretest is fetched
	tester.AssertNextCallIs(deps.fetchPretestCommand).InjectReturns(pretestCommand)
	// Then the pretest command is run
	tester.AssertNextCallIs(deps.runSubprocess, pretestCommand).InjectReturns(nil)
	// Then the function is done
	tester.AssertDoneWithin(time.Second)
	// Then the function passed
	tester.AssertReturned(result)
}

func FuzzPretestHappyPath(f *testing.F) {
	f.Fuzz(rapid.MakeFuzz(rapidPretestHappyPath))
}

func TestPretestHappyPath(t *testing.T) {
	t.Parallel()
	rapid.Check(t, rapidPretestHappyPath)
}

func rapidPretestSubprocessFail(rapidTester *rapid.T) {
	// Given test needs
	relay := protest.NewCallRelay()
	tester := &protest.RelayTester{T: rapidTester, Relay: relay}
	// Given inputs
	deps := newPretestDeps(relay)
	pretestCommand := rapid.SliceOf(rapid.String()).Draw(rapidTester, "pretestCommand")
	// this is explicitly a dynamic error for testing
	subprocessError := errors.New(rapid.String().Draw(rapidTester, "subprocessError")) //nolint:goerr113
	// Given outputs
	result := false

	// When the func is run
	tester.Start(pretest, deps)

	// Then the pretest is fetched
	tester.AssertNextCallIs(deps.fetchPretestCommand).InjectReturns(pretestCommand)
	// Then the pretest command is run
	tester.AssertNextCallIs(deps.runSubprocess, pretestCommand).InjectReturns(subprocessError)
	// Then the function is done
	tester.AssertDoneWithin(time.Second)
	// Then the function failed
	tester.AssertReturned(result)
}

func FuzzPretestSubprocessFail(f *testing.F) {
	f.Fuzz(rapid.MakeFuzz(rapidPretestSubprocessFail))
}

func TestPretestSubprocessFail(t *testing.T) {
	t.Parallel()
	rapid.Check(t, rapidPretestSubprocessFail)
}
