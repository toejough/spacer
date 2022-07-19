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
			mutate := func() bool {
				actualCalls.Push("mutate")
				return mReturn
			}
			reportT := func(r bool) {
				actualCalls.Push("report")
				reportArgs = r
			}
			exitT := func(r bool) {
				actualCalls.Push("exit")
				exitArgs = r
			}

			// When run is called
			run{mutate, reportT, exitT}.f()

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
			iReturn := iterator{}
			// Given test vars
			var (
				mutatorArgs iterator
			)
			// Given dependencies
			newIteratorT := func() iterator {
				actualCalls.Push("new file iterator")
				return iReturn
			}
			recursiveMutatorT := func(iterator) bool {
				actualCalls.Push("recursive mutator")
				return rmReturn
			}

			// When called
			mReturn := mutate{newIteratorT, recursiveMutatorT}.f()

			// Then get the file iterator
			protest.RequireCall(t, "new file iterator", actualCalls.MustPop(t))
			// Then call the recursive mutator with the iterator
			protest.RequireCall(t, "recursive mutator", actualCalls.MustPop(t))
			protest.RequireArgs(t, iReturn, mutatorArgs, diffIterator)
			// And no more calls are made
			protest.RequireEmpty(t, actualCalls)
			// And the output of the mutator is returned
			protest.RequireReturn(t, mReturn, rmReturn, diffBool)
		})
	}
}

/*
Properties
* overall flow is:
  * for each go file under pwd, for each pattern, for each instance of this pattern in this file, serially
  * report instance, mutate, test with CLI command (capture but don't act on result yet), restore, report
  * stop testing after first error/failure
  * report overall candidate/error/mutant status
  * exit with
    * 1 if no candidates were found
    * 1 if any errors occurred
    * 1 if any mutants survived
    * 0 if there were mutants and all of them were caught
*/
