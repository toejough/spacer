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
	exit                         struct{ code returnCodes }
	testMutationTypes            struct{ returnOneShot *protest.FIFO[mutationResult] }
	verifyTestsPassWithNoMutants struct{ returnOneShot *protest.FIFO[bool] }
)

func (rdm *runDepsMock) verifyTestsPassWithNoMutants() bool {
	returnOneShot := protest.NewOneShotFIFO[bool]("verifyTestsPassWithNoMutantsReturn")

	rdm.calls.Push(verifyTestsPassWithNoMutants{returnOneShot: returnOneShot})

	return returnOneShot.MustPop(rdm.t)
}

func (rdm *runDepsMock) testMutationTypes() mutationResult {
	returnOneShot := protest.NewOneShotFIFO[mutationResult]("testMutationTypesReturn")

	rdm.calls.Push(testMutationTypes{returnOneShot: returnOneShot})

	return returnOneShot.MustPop(rdm.t)
}

func (rdm *runDepsMock) exit(code returnCodes) {
	rdm.calls.Push(exit{code: code})
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

	// The mutant catcher is tested
	verifyCall := new(verifyTestsPassWithNoMutants)
	deps.calls.MustPopAs(t, verifyCall)

	// When the mutant catcher returns true
	verifyCall.returnOneShot.Push(true)

	// Then mutation type testing is done
	mutationTypesCall := new(testMutationTypes)
	deps.calls.MustPopAs(t, mutationTypesCall)

	// When the testing returns all caught
	mutationTypesCall.returnOneShot.Push(mutationResult{result: experimentResultAllCaught, err: nil})

	// Then the program exits
	deps.calls.MustPopEqualTo(t, exit{code: returnCodePass})
	// and there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
}

func TestRunTestsFailWithoutMutants(t *testing.T) {
	t.Parallel()

	deps := newMockedDeps(t)

	// When the func is run
	go func() {
		run(deps)
		deps.close()
	}()

	// The mutant catcher is tested
	verifyCall := new(verifyTestsPassWithNoMutants)
	deps.calls.MustPopAs(t, verifyCall)

	// When the mutant catcher returns true
	verifyCall.returnOneShot.Push(false)

	// Then the program exits
	deps.calls.MustPopEqualTo(t, exit{code: returnCodeTestsFailWithNoMutations})
	// and there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
}
