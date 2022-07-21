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
# Basic flow
* for each candidate
* test serially
* stop testing after error or failure
* exit with
    * 1 if no candidates were found
    * 1 if any errors occurred
    * 1 if any mutants survived
    * 0 if there were mutants and all of them were caught

# Candidate identification
* for each go file under the PWD
* for each pattern
* add all instances of the pattern

# testing
* candidate mutated
* test command from CLI run
* file restored
* test result returned

# Properties
* candidate mutated
    * given a candidate
    * when the test happens
    * then the candidate mutation is written to the file
* test command run
    * given the test command
    * when the test happens
    * then the test command is run
* file restored
    * given a candidate
    * when the test happens
    * the file is restored after the test command returns
* test command result returned
    * given a test command result
    * the test returns it
* file identification restrained to PWD
    * given a lot of files
    * when file scanning happens
    * only the files below the PWD are returned
* file types restricted to go
    * given a lot of files
    * when file scanning happens
    * only the go files are returned
* all patterns searched
    * given files and patterns
    * every file is searched for every pattern
* all results returned
    * given files and patterns
    * every result is returned from the searches
* all instances tested without early exit
    * given instances of patterns across files
    * and no failures/errors from testing
    * when run
    * then all instances are tested
* tests run serially
    * given instances of patterns across files
    * when run
    * then all tests run are run serially
* testing stops on error
    * given instances of patterns across files
    * and errors from testing
    * when run
    * then testing stops after the first error
* testing stops on failure
    * given instances of patterns across files
    * and failures from testing
    * when run
    * then testing stops after the first failure
* exit with 1 if no candidates found
    * given no instances
    * when run
    * exits with 1
* exit with 1 if errors found
    * given instances of patterns across files
    * and errors from testing
    * when run
    * exits with 1
* exit with 1 if failures found
    * given instances of patterns across files
    * and failures from testing
    * when run
    * exits with 1
* exit with 0 if all mutants caught
    * given instances of patterns across files
    * and no failures or errors from testing
    * when run
    * exits with 0
*/

// OR
/*
For every go file under PWD...
    For every pattern in hardcoded list, search...
        For every match...
            Mutate the file...
            Run the command from the CLI...
            Restore the file...
            If there was an error, exit with 1...
            If there the test passed, exit with 1...
If there were no matches found, exit with 1...
Exit with 0...

TestEveryGoFileUnderPWD
    Given go & non-go files under and not under PWD
    Given failing test results
    When run
    Then All go files under PWD are searched
    Then no non-go files are searched
    Then no files not under PWD are searched
TestSearchEveryPatternInList
    Given files
    Given patterns
    When run
    Then every file is searched for every pattern
    Then no other searches are carried out
TestEveryMatch
    Given files
    Given patterns
    Given matches for those patterns in those files
    When run
    Then every match is tested
TestTestPass
    Given a passing instance
    Given a CLI command
    When run
    Then mutation is performed
    Then file is saved
    Then CLI command is run
    Then file is restored
    Then pass is returned
TestTestFail
    Given a failing instance
    Given a CLI command
    When run
    Then mutation is performed
    Then file is saved
    Then CLI command is run
    Then file is restored
    Then fail is returned
TestTestError
    Given a Erroring instance
    Given a CLI command
    When run
    Then mutation is performed
    Then file is saved
    Then CLI command is run
    Then file is restored
    Then Error is returned
TestMutateRunRestoreExitError
    Given files
    Given patterns
    Given matches for those patterns in those files
    Given at least one error
    When run
    Then every match is tested up to and including the error case
    Then no other matches are tested
    Then exit with 1
TestMutateRunRestoreExitPass
    Given files
    Given patterns
    Given matches for those patterns in those files
    Given at least one pass
    When run
    Then every match is tested up to and including the passing case
    Then no other matches are tested
    Then exit with 1
TestNoMatchesExit
    Given files
    Given patterns
    Given no matches for those patterns in those files
    When run
    Then exit with 1
TestAllCaughtExit
    Given files
    Given patterns
    Given matches for those patterns in those files
    Given all failures
    When run
    Then exit with 0
*/
