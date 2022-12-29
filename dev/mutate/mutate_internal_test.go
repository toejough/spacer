package main

import (
	"fmt"
	"spacer/dev/protest"
	"testing"

	"pgregory.net/rapid"
)

type (
	runDepsMock struct {
		calls *protest.FIFO[interface{}]
		t     tester
	}
	exit                         struct{ code returnCodes }
	testMutationTypes            struct{ returnOneShot *protest.FIFO[mutationResult] }
	verifyTestsPassWithNoMutants struct{ returnOneShot *protest.FIFO[error] }
	tester                       interface {
		Helper()
		Fatal(...any)
	}
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

	rapid.Check(t, func(test *rapid.T) {
		deps := newMockedDeps(test)

		// When the func is run
		go func() {
			run(deps)
			deps.close()
		}()

		// The mutant catcher is tested
		verifyCall := new(verifyTestsPassWithNoMutants)
		deps.calls.MustPopAs(test, verifyCall)

		// When the mutant catcher returns no error
		// Don't grouse about the dynammic error here - it's the whole point
		verifyCall.returnOneShot.Push(fmt.Errorf(rapid.String().Draw(test, "mutationTypesError"))) //nolint: goerr113

		// Then the program exits
		deps.calls.MustPopEqualTo(test, exit{code: returnCodeTestsFailWithNoMutations})
		// and there are no more dependency calls
		deps.calls.MustConfirmClosed(test)
	})
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

	rapid.Check(t, func(test *rapid.T) {
		deps := newMockedDeps(test)

		// When the func is run
		go func() {
			run(deps)
			deps.close()
		}()

		// The mutant catcher is tested
		verifyCall := new(verifyTestsPassWithNoMutants)
		deps.calls.MustPopAs(test, verifyCall)

		// When the mutant catcher returns no error
		verifyCall.returnOneShot.Push(nil)

		// Then mutation type testing is done
		mutationTypesCall := new(testMutationTypes)
		deps.calls.MustPopAs(test, mutationTypesCall)

		// When the testing returns all caught
		mutationTypesCall.returnOneShot.Push(mutationResult{
			result: experimentResultUndetectedMutants,
			// Don't grouse about the dynammic error here - it's the whole point
			err: fmt.Errorf(rapid.String().Draw(test, "mutationTypesError")), //nolint: goerr113
		})

		// Then the program exits
		deps.calls.MustPopEqualTo(test, exit{code: returnCodeFail})
		// and there are no more dependency calls
		deps.calls.MustConfirmClosed(test)
	})
}
