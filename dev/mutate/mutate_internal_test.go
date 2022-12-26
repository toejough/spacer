package main

import (
	"spacer/dev/protest"
	"testing"
)

// TODO make deps an interface
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
	calls := protest.NewFIFO[string]("calls")
	exitArgs := protest.NewFIFO[returnCodes]("exitArgs")
	verifyMutantCatcherPassesReturns := protest.NewFIFO[bool]("verifyMutantCatcherPassesReturns")
	testMutationTypesReturns := protest.NewFIFO[mutationResult]("testMutationTypesReturns")

	return mockRunDeps{
		calls:                            calls,
		exitArgs:                         exitArgs,
		verifyMutantCatcherPassesReturns: verifyMutantCatcherPassesReturns,
		testMutationTypesReturns:         testMutationTypesReturns,
		deps: runDeps{
			announceMutationTesting: func() { calls.Push("announceMutationTesting") },
			verifyMutantCatcherPasses: func() bool {
				calls.Push("verifyMutantCatcherPasses")
				return verifyMutantCatcherPassesReturns.MustGetNext(t)
			},
			testMutationTypes: func() mutationResult {
				calls.Push("testMutationTypes")
				return testMutationTypesReturns.MustGetNext(t)
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
	// TODO use enums instead of strings for function names
	deps.calls.RequireNext(t, "announceMutationTesting")
	// And the mutant catcher is tested
	deps.calls.RequireNext(t, "verifyMutantCatcherPasses")

	// When the mutant catcher returns true
	deps.verifyMutantCatcherPassesReturns.Push(true)

	// Then mutation type testing is done
	deps.calls.RequireNext(t, "testMutationTypes")

	// When the testing is all caught
	deps.testMutationTypesReturns.Push(mutationResult{result: experimentResultAllCaught, err: nil})

	// Then the program exits
	deps.calls.RequireNext(t, "exit")
	// and does so with a passing %return code
	deps.exitArgs.RequireNext(t, returnCodePass)
	// and there are no more dependency calls
	{
		// TODO make this follow the t/no-t example of Get. CheckFinal vs RequireFinal?
		// TODO split into two? CheckClosed, CheckDrained? A helper for both? what does everyone else do about channels?
		err := deps.calls.RequireClosedAndEmpty()
		if err != nil {
			t.Fatal(err)
		}
	}
}
