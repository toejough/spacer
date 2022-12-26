package main

import (
	"spacer/dev/protest"
	"testing"
	"time"

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
	})
	exitArgs := protest.NewFIFO("exitArgs", protest.FIFODeps[returnCodes]{Differ: returnCodeDiff})
	verifyMutantCatcherPassesReturns := protest.NewFIFO("verifyMutantCatcherPassesReturns", protest.FIFODeps[bool]{
		Differ: boolDiff,
	})
	testMutationTypesReturns := protest.NewFIFO("testMutationTypesReturns", protest.FIFODeps[mutationResult]{
		Differ: mutationResultDiff,
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
				toReturn, err := verifyMutantCatcherPassesReturns.WaitForNext(1 * time.Second)
				if err != nil {
					t.Fatal(err)
				}

				return toReturn
			},
			testMutationTypes: func() mutationResult {
				calls.Push("testMutationTypes")
				toReturn, err := testMutationTypesReturns.WaitForNext(1 * time.Second)
				if err != nil {
					t.Fatal(err)
				}

				return toReturn
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
	{
		called, err := deps.calls.WaitForNext(1 * time.Second)
		if err != nil {
			t.Fatal(err)
		}

		expected := "announceMutationTesting"
		if called != expected {
			t.Fatalf("expected %s but %s was called instead", expected, called)
		}
	}
	// And the mutant catcher is tested
	{
		called, err := deps.calls.WaitForNext(1 * time.Second)
		if err != nil {
			t.Fatal(err)
		}

		expected := "verifyMutantCatcherPasses"
		if called != expected {
			t.Fatalf("expected %s but %s was called instead", expected, called)
		}
	}

	// When the mutant catcher returns true
	deps.verifyMutantCatcherPassesReturns.Push(true)

	// Then mutation type testing is done
	{
		called, err := deps.calls.WaitForNext(1 * time.Second)
		if err != nil {
			t.Fatal(err)
		}

		expected := "testMutationTypes"
		if called != expected {
			t.Fatalf("expected %s but %s was called instead", expected, called)
		}
	}

	// When the testing is all caught
	deps.testMutationTypesReturns.Push(mutationResult{result: experimentResultAllCaught, err: nil})

	// Then the program exits
	{
		called, err := deps.calls.WaitForNext(1 * time.Second)
		if err != nil {
			t.Fatal(err)
		}

		expected := "exit"
		if called != expected {
			t.Fatalf("expected %s but %s was called instead", expected, called)
		}
	}
	// and does so with a passing %return code
	{
		returned, err := deps.exitArgs.WaitForNext(1 * time.Second)
		if err != nil {
			t.Fatal(err)
		}

		expected := returnCodePass
		if returned != expected {
			t.Fatalf("expected %v but %v was called instead", expected, returned)
		}
	}
	// and there are no more dependency calls
	{
		err := deps.calls.RequireClosedAndEmpty()
		if err != nil {
			t.Fatal(err)
		}
	}
}
