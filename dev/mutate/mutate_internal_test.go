package main

// The main idea for all the unit tests is to test the behavior we care about _at this level_.
// This means we validate the calls to dependencies _at this level_ (critically, _not_ subdependency calls).
// Leave "and now xyz is happening" testing to the thing that is making it happen.
// For example, for "run", we do _not_ care where the test command is coming from, how it is run, or how its output is
// conveyed.

import (
	"spacer/dev/protest"
	"testing"
)

type (
	runDepsMock struct {
		calls *protest.FIFO[interface{}]
		t     tester
	}
	announceStartingCall             struct{}
	verifyTestsPassWithNoMutantsCall struct {
		returnOneShot *protest.FIFO[bool]
	}
	testMutationsCall struct {
		returnOneShot *protest.FIFO[bool]
	}
	announceEndingCall struct{}
	exitArgs           struct{ passed bool }
	exitCall           struct {
		args exitArgs
	}
	tester interface {
		Helper()
		Fatal(...any)
	}
)

func (rdm *runDepsMock) announceStarting() {
	rdm.calls.Push(announceStartingCall{})
}

func (rdm *runDepsMock) verifyTestsPassWithNoMutants() bool {
	returnOneShot := protest.NewOneShotFIFO[bool]("verifyTestsPassWithNoMutantsReturn")

	rdm.calls.Push(verifyTestsPassWithNoMutantsCall{returnOneShot: returnOneShot})

	return returnOneShot.MustPop(rdm.t)
}

func (rdm *runDepsMock) testMutations() bool {
	returnOneShot := protest.NewOneShotFIFO[bool]("testMutationsReturn")

	rdm.calls.Push(testMutationsCall{returnOneShot: returnOneShot})

	return returnOneShot.MustPop(rdm.t)
}

func (rdm *runDepsMock) announceEnding() {
	rdm.calls.Push(announceEndingCall{})
}

func (rdm *runDepsMock) exit(passed bool) {
	rdm.calls.Push(exitCall{args: exitArgs{passed: passed}})
}

func (rdm *runDepsMock) close() {
	rdm.calls.Close()
}

func newMockedDeps(t tester) *runDepsMock {
	t.Helper()

	return &runDepsMock{
		calls: protest.NewFIFO[any]("calls"),
		t:     t,
	}
}

func TestRunHappyPath(t *testing.T) {
	t.Parallel()

	deps := newMockedDeps(t)

	// When the func is run
	go func() {
		run(deps)
		deps.close()
	}()

	// Then the program announces itself
	deps.calls.MustPopEqualTo(t, announceStartingCall{})
	// And the tester is run to ensure it passes prior to applying mutations
	verifyCall := new(verifyTestsPassWithNoMutantsCall)
	deps.calls.MustPopAs(t, verifyCall)

	// When the tester passes
	verifyCall.returnOneShot.Push(true)

	// Then mutation type testing is performed
	testMutationsCall := new(testMutationsCall)
	deps.calls.MustPopAs(t, testMutationsCall)

	// When the testing returns all caught
	testMutationsCall.returnOneShot.Push(true)

	// Then program announces it's exiting
	deps.calls.MustPopEqualTo(t, announceEndingCall{})
	// And the program exits
	deps.calls.MustPopEqualTo(t, exitCall{args: exitArgs{passed: true}})
	// and there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
}

func TestRunTesterFailsBeforeAnyMutations(t *testing.T) {
	t.Parallel()

	deps := newMockedDeps(t)

	// When the func is run
	go func() {
		run(deps)
		deps.close()
	}()

	// Then the program announces itself
	deps.calls.MustPopEqualTo(t, announceStartingCall{})
	// And the tester is run to ensure it passes prior to applying mutations
	verifyCall := new(verifyTestsPassWithNoMutantsCall)
	deps.calls.MustPopAs(t, verifyCall)

	// When the tester passes
	verifyCall.returnOneShot.Push(false)

	// Then program announces it's exiting
	deps.calls.MustPopEqualTo(t, announceEndingCall{})
	// And the program exits
	deps.calls.MustPopEqualTo(t, exitCall{args: exitArgs{passed: false}})
	// and there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
}

func TestRunMutationTestsFail(t *testing.T) {
	t.Parallel()

	deps := newMockedDeps(t)

	// When the func is run
	go func() {
		run(deps)
		deps.close()
	}()

	// Then the program announces itself
	deps.calls.MustPopEqualTo(t, announceStartingCall{})
	// And the tester is run to ensure it passes prior to applying mutations
	verifyCall := new(verifyTestsPassWithNoMutantsCall)
	deps.calls.MustPopAs(t, verifyCall)

	// When the tester passes
	verifyCall.returnOneShot.Push(true)

	// Then mutation type testing is performed
	testMutationsCall := new(testMutationsCall)
	deps.calls.MustPopAs(t, testMutationsCall)

	// When the testing returns all caught
	testMutationsCall.returnOneShot.Push(false)

	// Then program announces it's exiting
	deps.calls.MustPopEqualTo(t, announceEndingCall{})
	// And the program exits
	deps.calls.MustPopEqualTo(t, exitCall{args: exitArgs{passed: false}})
	// and there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
}
