package main

import (
	"spacer/dev/protest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func stringDiff(e, a string) string {
	return cmp.Diff(e, a)
}

func boolDiff(e, a bool) string {
	return cmp.Diff(e, a)
}

func mutationResultDiff(e, a mutationResult) string {
	return cmp.Diff(e, a)
}

func returnCodeDiff(e, a returnCodes) string {
	return cmp.Diff(e, a)
}

// TODO make deps an interface
// TODO remove t requirement from protest inits
// TODO move protest closes into the mock deps implementation

type mockRunDeps struct {
	deps                             runDeps
	calls                            *protest.FIFO[string]
	exitArgs                         *protest.FIFO[returnCodes]
	verifyMutantCatcherPassesReturns *protest.FIFO[bool]
	testMutationTypesReturns         *protest.FIFO[mutationResult]
}

func newMockedDeps(t *testing.T) mockRunDeps {
	t.Helper()

	// Given Call/Arg/Return FIFOS
	calls := protest.NewFIFO("calls", protest.FIFODeps[string]{
		Differ: stringDiff,
		T:      t,
	})
	exitArgs := protest.NewFIFO("exitArgs", protest.FIFODeps[returnCodes]{Differ: returnCodeDiff, T: t})
	verifyMutantCatcherPassesReturns := protest.NewFIFO("verifyMutantCatcherPassesReturns", protest.FIFODeps[bool]{
		Differ: boolDiff,
		T:      t,
	})
	testMutationTypesReturns := protest.NewFIFO("testMutationTypesReturns", protest.FIFODeps[mutationResult]{
		Differ: mutationResultDiff,
		T:      t,
	})

	return mockRunDeps{
		calls:                            calls,
		exitArgs:                         exitArgs,
		verifyMutantCatcherPassesReturns: verifyMutantCatcherPassesReturns,
		testMutationTypesReturns:         testMutationTypesReturns,
		deps: runDeps{
			announceMutationTesting: func() { calls.Push("announceMutationTesting") },
			verifyMutantCatcherPasses: func() bool {
				calls.Push("verifyMutantCatcherPasses")
				return verifyMutantCatcherPassesReturns.GetNext()
			},
			testMutationTypes: func() mutationResult {
				calls.Push("testMutationTypes")
				return testMutationTypesReturns.GetNext()
			},
			exit: func(code returnCodes) {
				calls.Push("exit")
				exitArgs.Push(code)
			},
		},
	}
}

func TestRunHappyPath(t *testing.T) {
	t.Parallel()

	deps := newMockedDeps(t)

	// When the func is run
	go func() {
		run(deps.deps)
		deps.calls.Close()
	}()

	// Then mutation testing is announced
	deps.calls.RequireNext("announceMutationTesting")
	// And the mutant catcher is tested
	deps.calls.RequireNext("verifyMutantCatcherPasses")

	// When the mutant catcher returns true
	deps.verifyMutantCatcherPassesReturns.Push(true)

	// Then mutation type testing is done
	deps.calls.RequireNext("testMutationTypes")

	// When the testing is all caught
	deps.testMutationTypesReturns.Push(mutationResult{result: experimentResultAllCaught, err: nil})

	// Then the program exits
	deps.calls.RequireNext("exit")
	// and does so with a passing return code
	deps.exitArgs.RequireNext(returnCodePass)
	// and there are no more dependency calls
	deps.calls.RequireClosedAndEmpty()
}
