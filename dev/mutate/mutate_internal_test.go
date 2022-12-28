package main

import (
	"spacer/dev/protest"
	"testing"
)

type (
	runDepsMock struct {
		calls *protest.FIFO[interface{}]
		t     *testing.T
	}
	exitMock                      struct{ code returnCodes }
	testMutationTypesMock         struct{ returnOneShot *protest.FIFO[mutationResult] }
	verifyMutantCatcherPassesMock struct{ returnOneShot *protest.FIFO[bool] }
	announceMutationTestingMock   struct{}
)

func (rdm *runDepsMock) announceMutationTesting() {
	rdm.calls.Push(announceMutationTestingMock{})
}

func (rdm *runDepsMock) verifyMutantCatcherPasses() bool {
	returnOneShot := protest.NewOneShotFIFO[bool]("verifyMutantCatcherPassesReturns")

	rdm.calls.Push(verifyMutantCatcherPassesMock{returnOneShot: returnOneShot})

	return returnOneShot.MustPop(rdm.t)
}

func (rdm *runDepsMock) testMutationTypes() mutationResult {
	returnOneShot := protest.NewOneShotFIFO[mutationResult]("testMutationTypesReturns")

	rdm.calls.Push(testMutationTypesMock{returnOneShot: returnOneShot})

	return returnOneShot.MustPop(rdm.t)
}

func (rdm *runDepsMock) exit(code returnCodes) {
	rdm.calls.Push(exitMock{code: code})
}

func (rdm *runDepsMock) close() {
	rdm.calls.Close()
}

func newMockedDeps(t *testing.T) *runDepsMock {
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

	// Then mutation testing is announced
	deps.calls.MustPopEqualTo(t, announceMutationTestingMock{})
	// And the mutant catcher is tested
	verifyCall := new(verifyMutantCatcherPassesMock)
	deps.calls.MustPopAs(t, verifyCall)

	// When the mutant catcher returns true
	verifyCall.returnOneShot.Push(true)

	// Then mutation type testing is done
	mutationTypesCall := new(testMutationTypesMock)
	deps.calls.MustPopAs(t, mutationTypesCall)

	// When the testing returns all caught
	mutationTypesCall.returnOneShot.Push(mutationResult{result: experimentResultAllCaught, err: nil})

	// Then the program exits
	deps.calls.MustPopEqualTo(t, exitMock{code: returnCodePass})
	// and there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
}
