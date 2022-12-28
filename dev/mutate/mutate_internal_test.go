package main

import (
	"spacer/dev/protest"
	"testing"
)

// TODO make deps an interface

type mockRunDeps struct {
	deps  runDeps
	calls *protest.FIFO[interface{}]
}

func (m *mockRunDeps) close() {
	m.calls.Close()
}

func newMockedDeps(t *testing.T) mockRunDeps {
	t.Helper()

	// Given Call/Arg/Return FIFOS
	calls := protest.NewFIFO[any]("calls")

	return mockRunDeps{
		calls: calls,
		deps: runDeps{
			announceMutationTesting: func() { calls.Push(announceMutationTestingMock{}) },
			verifyMutantCatcherPasses: func() bool {
				returnOneShot := protest.NewOneShotFIFO[bool]("verifyMutantCatcherPassesReturns")
				calls.Push(verifyMutantCatcherPassesMock{returnFifo: returnOneShot})

				return returnOneShot.MustPop(t)
			},
			testMutationTypes: func() mutationResult {
				returnOneShot := protest.NewOneShotFIFO[mutationResult]("testMutationTypesReturns")
				calls.Push(testMutationTypesMock{returnFifo: returnOneShot})

				return returnOneShot.MustPop(t)
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
		deps.close()
	}()

	// Then mutation testing is announced
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
	err := deps.calls.ConfirmClosed()
	if err != nil {
		t.Fatal(err)
	}
}
