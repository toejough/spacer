package main

import (
	"spacer/dev/protest"
	"testing"
)

// TODO make deps an interface
// TODO move protest closes into the mock deps implementation

type mockRunDeps struct {
	deps                             runDeps
	calls                            *protest.FIFO[interface{}]
	exitArgs                         *protest.FIFO[returnCodes]
	verifyMutantCatcherPassesReturns *protest.FIFO[bool]
	testMutationTypesReturns         *protest.FIFO[mutationResult]
}

func newMockedDeps(t *testing.T) mockRunDeps {
	t.Helper()

	// Given Call/Arg/Return FIFOS
	calls := protest.NewFIFO[any]("calls")
	exitArgs := protest.NewFIFO[returnCodes]("exitArgs")
	verifyMutantCatcherPassesReturns := protest.NewFIFO[bool]("verifyMutantCatcherPassesReturns")
	testMutationTypesReturns := protest.NewFIFO[mutationResult]("testMutationTypesReturns")

	return mockRunDeps{
		calls:                            calls,
		exitArgs:                         exitArgs,
		verifyMutantCatcherPassesReturns: verifyMutantCatcherPassesReturns,
		testMutationTypesReturns:         testMutationTypesReturns,
		deps: runDeps{
			announceMutationTesting: func() { calls.Push(announceMutationTestingMock{}) },
			verifyMutantCatcherPasses: func() bool {
				calls.Push(verifyMutantCatcherPassesMock{returnFifo: verifyMutantCatcherPassesReturns})
				return verifyMutantCatcherPassesReturns.MustPop(t)
			},
			testMutationTypes: func() mutationResult {
				calls.Push(testMutationTypesMock{returnFifo: testMutationTypesReturns})
				return testMutationTypesReturns.MustPop(t)
			},
			exit: func(code returnCodes) {
				calls.Push(exitMock{code: code})
			},
		},
	}
}

type exitMock struct {
	code returnCodes
}

type testMutationTypesMock struct {
	returnFifo *protest.FIFO[mutationResult]
}

type verifyMutantCatcherPassesMock struct {
	returnFifo *protest.FIFO[bool]
}

type announceMutationTestingMock struct{}

func TestRunHappyPath(t *testing.T) {
	t.Parallel()

	deps := newMockedDeps(t)

	// When the func is run
	go func() {
		run(deps.deps)
		deps.calls.Close()
	}()

	// Then mutation testing is announced
	// TODO use a one-shot mechanism for the return values instead of FIFO's
	deps.calls.MustPopEqualTo(t, announceMutationTestingMock{})
	// And the mutant catcher is tested
	verifyCall := new(verifyMutantCatcherPassesMock)
	deps.calls.MustPopAs(t, verifyCall)
	// When the mutant catcher returns true
	verifyCall.returnFifo.Push(true)

	// Then mutation type testing is done
	mutationTypesCall := new(testMutationTypesMock)
	deps.calls.MustPopAs(t, mutationTypesCall)
	// When the testing is all caught
	mutationTypesCall.returnFifo.Push(mutationResult{result: experimentResultAllCaught, err: nil})

	// Then the program exits
	deps.calls.MustPopEqualTo(t, exitMock{code: returnCodePass})
	// and there are no more dependency calls
	// TODO make this follow the t/no-t example of Get. GetClosed & GetClosedWithin?
	// TODO need a within variant
	err := deps.calls.RequireClosedAndEmpty()
	if err != nil {
		t.Fatal(err)
	}
}
