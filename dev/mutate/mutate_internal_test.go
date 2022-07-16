package main

import (
	"fmt"
	"spacer/dev/protest"
	"testing"
)

func diffBool(expected, actual bool) string {
	if expected != actual {
		return fmt.Sprintf("%t != %t", expected, actual)
	}

	return ""
}

func TestRun(t *testing.T) {
	t.Parallel()

	for _, mutationReturn := range []bool{true, false} { //nolint:paralleltest
		// linter can't tell we're using the range value in the test
		// Given a return from the mutation func
		mReturn := mutationReturn

		t.Run(fmt.Sprintf("%t", mutationReturn), func(t *testing.T) {
			t.Parallel()

			// Given test objects
			actualCalls := protest.FIFO[string]{}
			// Given test vars
			var reportArgs bool
			var exitArgs bool
			// Given dependencies
			mutateT := func() bool {
				actualCalls.Push("mutate")

				return mReturn
			}
			report := func(r bool) {
				actualCalls.Push("report")
				reportArgs = r
			}
			exit := func(r bool) {
				actualCalls.Push("exit")
				exitArgs = r
			}

			// When run is called
			run{mutateT, report, exit}.f()

			// Then mutate is called
			protest.RequireCall(t, "mutate", actualCalls.MustPop(t))
			// Then report is called with mutate's output
			protest.RequireCall(t, "report", actualCalls.MustPop(t))
			protest.RequireArgs(t, mReturn, reportArgs, diffBool)
			// Then exit is called with mutate's output
			protest.RequireCall(t, "exit", actualCalls.MustPop(t))
			protest.RequireArgs(t, mReturn, exitArgs, diffBool)
			// And no more calls are made
			protest.RequireEmpty(t, actualCalls)
		})
	}
}

func diffIterator(expected, actual iterator) string { return "" }

type iterator struct{}

func TestMutate(t *testing.T) {
	t.Parallel()

	// Given test objects
	actualCalls := protest.FIFO[string]{}
	// Given test vars
	var (
		mutatorArgs iterator
		iReturn     iterator
		mrReturn    bool
	)
	// Given dependencies

	// When called
	mReturn := mutate()

	// Then get the file iterator
	protest.RequireCall(t, "new file iterator", actualCalls.MustPop(t))
	// Then call the recursive mutator with the iterator
	protest.RequireCall(t, "recursive mutator", actualCalls.MustPop(t))
	protest.RequireArgs(t, iReturn, mutatorArgs, diffIterator)
	// And no more calls are made
	protest.RequireEmpty(t, actualCalls)
	// And the output of the mutator is returned
	protest.RequireReturn(t, mReturn, mrReturn, diffBool)
}
