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
		deps  *runDeps
	}
	verifyTestsPassWithNoMutantsCall protest.CallWithNoArgs[bool]
	testMutationsCall                protest.CallWithNoArgs[bool]
	tester                           interface {
		Helper()
		Fatal(...any)
	}
)

func (rdm *runDepsMock) close() {
	rdm.calls.Close()
}

func newMockedDeps(test tester) *runDepsMock {
	test.Helper()

	calls := protest.NewFIFO[any]("calls")

	return &runDepsMock{
		calls: calls,
		t:     test,
		deps: &runDeps{
			pretest: func() bool {
				return protest.ManageCallWithNoArgs[verifyTestsPassWithNoMutantsCall](test, calls)
			},
			testMutations: func() bool { return protest.ManageCallWithNoArgs[testMutationsCall](test, calls) },
		},
	}
}

func TestRunHappyPath(t *testing.T) {
	t.Parallel()

	// Given inputs & outputs
	deps := newMockedDeps(t)

	var passes bool

	// When the func is run
	go func() {
		passes = run(deps.deps)
		deps.close()
	}()

	// Then the tester is run to ensure it passes prior to applying mutations
	verifyCall := new(verifyTestsPassWithNoMutantsCall)
	deps.calls.MustPopAs(t, verifyCall)

	// When the tester passes
	verifyCall.ReturnOneShot.Push(true)

	// Then mutation type testing is performed
	testMutationsCall := new(testMutationsCall)
	deps.calls.MustPopAs(t, testMutationsCall)

	// When the testing returns all caught
	testMutationsCall.ReturnOneShot.Push(true)

	// Then there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
	// and the return value is as expected
	protest.MustEqual(t, true, passes)
}

func TestRunTesterFailsBeforeAnyMutations(t *testing.T) {
	t.Parallel()

	deps := newMockedDeps(t)

	var passes bool

	// When the func is run
	go func() {
		passes = run(deps.deps)
		deps.close()
	}()

	// Then the tester is run to ensure it passes prior to applying mutations
	verifyCall := new(verifyTestsPassWithNoMutantsCall)
	deps.calls.MustPopAs(t, verifyCall)

	// When the tester fails
	verifyCall.ReturnOneShot.Push(false)

	// Then there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
	// and the return value is as expected
	protest.MustEqual(t, false, passes)
}

func TestRunMutationTestsFail(t *testing.T) {
	t.Parallel()

	deps := newMockedDeps(t)

	var passes bool

	// When the func is run
	go func() {
		passes = run(deps.deps)
		deps.close()
	}()

	// Then the tester is run to ensure it passes prior to applying mutations
	verifyCall := new(verifyTestsPassWithNoMutantsCall)
	deps.calls.MustPopAs(t, verifyCall)

	// When the tester passes
	verifyCall.ReturnOneShot.Push(true)

	// Then mutation type testing is performed
	testMutationsCall := new(testMutationsCall)
	deps.calls.MustPopAs(t, testMutationsCall)

	// When the testing returns all caught
	testMutationsCall.ReturnOneShot.Push(false)

	// Then there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
	// and the return value is as expected
	protest.MustEqual(t, false, passes)
}
