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
			mutate := func() (bool, error) {
				actualCalls.Push("mutate")
				// TODO: what about the error case?
				return mReturn, nil
			}
			// TODO: test the prod version of this
			reportT := func(r bool) {
				actualCalls.Push("report")
				reportArgs = r
			}
			// TODO: test the prod version of this
			exitT := func(r bool) {
				actualCalls.Push("exit")
				exitArgs = r
			}

			// When run is called
			run{mutate, reportT, exitT, nil}.f()

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
func diffErr(expected, actual error) string         { return "" }

func TestMutate(t *testing.T) {
	t.Parallel()

	for _, rMutationReturn := range []bool{true, false} { //nolint:paralleltest
		// linter can't tell we're using the range value in the test
		// Given a return from the mutation func
		rmReturn := rMutationReturn

		t.Run(fmt.Sprintf("%t", rMutationReturn), func(t *testing.T) {
			t.Parallel()
			// Given test objects
			actualCalls := protest.FIFO[string]{}
			iReturn := iterator{root: "test"}
			// Given test vars
			var (
				mutatorArgs iterator
			)
			// Given dependencies
			// TODO: test the prod version of this
			newIteratorT := func() (iterator, error) {
				actualCalls.Push("new file iterator")
				// TODO: what about the error case?
				return iReturn, nil
			}
			// TODO: test the prod version of this
			recursiveMutatorT := func(iterator) bool {
				actualCalls.Push("recursive mutator")
				return rmReturn
			}

			// When called
			mReturn, err := mutate{newIteratorT, recursiveMutatorT}.f()

			// Then get the file iterator
			protest.RequireCall(t, "new file iterator", actualCalls.MustPop(t))
			// Then call the recursive mutator with the iterator
			protest.RequireCall(t, "recursive mutator", actualCalls.MustPop(t))
			protest.RequireArgs(t, iReturn, mutatorArgs, diffIterator)
			// And no more calls are made
			protest.RequireEmpty(t, actualCalls)
			// And the output of the mutator is returned
			protest.RequireReturn(t, mReturn, rmReturn, diffBool)
			protest.RequireReturn(t, err, nil, diffErr)
		})
	}
}
