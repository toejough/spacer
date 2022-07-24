package main

import (
	"spacer/dev/protest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

/*
Run:
    Test files
    Report results
    If there was an error, exit with 1...
    If there the test passed, exit with 1...
    If there were no matches found, exit with 1...
    Exit with 0...
*/

func diffString(expected, actual string) string {
	return cmp.Diff(expected, actual)
}

func diffBool(expected, actual bool) string {
	return cmp.Diff(expected, actual)
}

func diffInt(expected, actual int) string {
	return cmp.Diff(expected, actual)
}

func TestRunAll(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		testResults bool
		exitCode    int
	}{
		"passing":          {true, 0},
		"uncaught mutant":  {false, 1},
		"error during run": {false, 1},
		"no matches found": {false, 1},
	}

	for name, testCase := range testTable { //nolint:paralleltest // we _are_ using the range value
		testCase := testCase // avoid using the loop-scoped value in the closure

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given call FIFO
			calls := protest.NewFIFO[string]("calls")
			exitArgs := protest.NewFIFO[int]("exitArgs")
			reportArgs := protest.NewFIFO[bool]("reportArgs")
			// Given test results
			testResults := testCase.testResults
			// Given expected exit code
			expectedExitCode := testCase.exitCode
			// When run
			run{
				testFiles: func() bool {
					calls.Push("testFiles")
					return testResults
				},
				reportResults: func(results bool) {
					calls.Push("reportResults")
					reportArgs.Push(results)
				},
				exit: func(code int) {
					calls.Push("exit")
					exitArgs.Push(code)
				},
			}.f()
			// Then testFiles is called and returns all passes
			protest.RequireNext(t, "testFiles", calls, diffString)
			// Then report called with results
			protest.RequireNext(t, "reportResults", calls, diffString)
			protest.RequireNext(t, testResults, reportArgs, diffBool)
			// Then exit called with 0
			protest.RequireNext(t, "exit", calls, diffString)
			protest.RequireNext(t, expectedExitCode, exitArgs, diffInt)
			// Then no more calls are made
			protest.RequireEmpty(t, calls)
		})
	}
}

/*
Test files:
    Get iterator for go files under PWD
    For every file:
        test all patterns
        early result return if error
        early result return if any verification failed
    return results

test all patterns:
    Get iterator for patterns
    for every patterns
        test all matches
        early result return if error
        early result return if any verification failed
    return results

test all matches:
    Get iterator for matches of pattern in file
    for every candidate found,
        test candidate
        early result return if error
        early result return if any verification failed
    return results

test candidate
    verify mutation caught
    Restore the file
    chain any error
    return result

verify mutation caught:
    Mutate the file...
    return error early
    Run the command from the CLI...
    return result
*/
