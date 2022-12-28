package main

import (
	"fmt"
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
	verifyTestsPassWithNoMutants struct{ returnOneShot *protest.FIFO[error] }
)

func (rdm *runDepsMock) verifyTestsPassWithNoMutants() error {
	returnOneShot := protest.NewOneShotFIFO[error]("verifyTestsPassWithNoMutantsReturn")

	rdm.calls.Push(verifyTestsPassWithNoMutants{returnOneShot: returnOneShot})

	// this is the specific error to return for the test
	return returnOneShot.MustPop(rdm.t) //nolint: wrapcheck
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

	// When the mutant catcher returns no error
	verifyCall.returnOneShot.Push(nil)

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

	// When the mutant catcher returns no error
	// TODO make this any error via rapid.
	// this is intentionally an arbitrary error.
	verifyCall.returnOneShot.Push(fmt.Errorf("arbitrary error")) //nolint: goerr113

	// Then the program exits
	deps.calls.MustPopEqualTo(t, exit{code: returnCodeTestsFailWithNoMutations})
	// and there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
}

func TestRunNoMutationCandidatesFound(t *testing.T) {
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

	// When the mutant catcher returns no error
	verifyCall.returnOneShot.Push(nil)

	// Then mutation type testing is done
	mutationTypesCall := new(testMutationTypes)
	deps.calls.MustPopAs(t, mutationTypesCall)

	// When the testing returns all caught
	mutationTypesCall.returnOneShot.Push(mutationResult{result: experimentResultNoCandidatesFound, err: nil})

	// Then the program exits
	deps.calls.MustPopEqualTo(t, exit{code: returnCodeNoCandidatesFound})
	// and there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
}

func TestRunUndetectedMutants(t *testing.T) {
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

	// When the mutant catcher returns no error
	verifyCall.returnOneShot.Push(nil)

	// Then mutation type testing is done
	mutationTypesCall := new(testMutationTypes)
	deps.calls.MustPopAs(t, mutationTypesCall)

	// When the testing returns all caught
	mutationTypesCall.returnOneShot.Push(mutationResult{result: experimentResultUndetectedMutants, err: nil})

	// Then the program exits
	deps.calls.MustPopEqualTo(t, exit{code: returnCodeFail})
	// and there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
}

func TestRunDetectionError(t *testing.T) {
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

	// When the mutant catcher returns no error
	verifyCall.returnOneShot.Push(nil)

	// Then mutation type testing is done
	mutationTypesCall := new(testMutationTypes)
	deps.calls.MustPopAs(t, mutationTypesCall)

	// When the testing returns all caught
	// TODO table drive this one - they should all be error, no matter the return type, and no matter the actual error
	// (use rapid)
	mutationTypesCall.returnOneShot.Push(mutationResult{
		result: experimentResultUndetectedMutants,
		// Don't grouse about the dynammic error here, it's supposed to be even _more_ dynamic (see above todo)
		err: fmt.Errorf("any arbitrary error"), //nolint: goerr113
	})

	// Then the program exits
	deps.calls.MustPopEqualTo(t, exit{code: returnCodeFail})
	// and there are no more dependency calls
	deps.calls.MustConfirmClosed(t)
}
