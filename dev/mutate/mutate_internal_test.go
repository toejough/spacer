package main

import (
	"fmt"
	"spacer/dev/protest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func stringDiff(e, a string) string {
	return cmp.Diff(e, a)
}

func returnCodeDiff(e, a returnCodes) string {
	return cmp.Diff(e, a)
}

func newMockedRunner(
	calls *protest.FIFO[string],
	exitCodes *protest.FIFO[returnCodes],
	mockedMutantCatcherResult bool,
	mockedMutationResult mutationResult,
) runner {
	return runner{
		announceMutationTesting: func() { calls.Push("announce mutation testing") },
		verifyMutantCatcherPasses: func() bool {
			calls.Push("verify mutant catcher passes prior to mutations")
            // calls.Push("verifyMutantCatcherPasses")
            // return verifyMutantCatcherPassesReturns.PopLeft()
			return mockedMutantCatcherResult
		},
		testMutationTypes: func() mutationResult {
			calls.Push("test mutation types")
			return mockedMutationResult
		},
		exit: func(code returnCodes) {
			calls.Push("exit")
			exitCodes.Push(code)
		},
	}
}

func TestRunHappyPath(t *testing.T) {
	t.Parallel()

	calls := protest.NewFIFO[string]("calls")
    // calls := protest.NewFIFO[string]("calls", protest.NewFIFODeps{
    //   t: t,
    //   differ: stringDiff,
    // })
	exitCodes := protest.NewFIFO[returnCodes]("exit codes")

	// Given happy path return values from dependencies
	mutantCatcherPasses := true
	allMutantsCaught := mutationResult{result: experimentResultAllCaught, err: nil}
	passCode := returnCodePass
	theRunner := newMockedRunner(calls, exitCodes, mutantCatcherPasses, allMutantsCaught)

	// When the func is run
	theRunner.run()

	// Then the program is announced
    // TODO: The comments here are smells that this call flow is not clean enough
    // calls.RequireNext("verifyMutantCatcherPasses")
    // verifyMutantCatcherPassesReturns.Push(true)
	protest.RequireNext(t, "announce mutation testing", calls, stringDiff)
	// And the mutant catcher is verified to pass prior to mutations
	protest.RequireNext(t, "verify mutant catcher passes prior to mutations", calls, stringDiff)
	// And the mutations are run
	protest.RequireNext(t, "test mutation types", calls, stringDiff)
	// And the program exits with 0
	protest.RequireNext(t, "exit", calls, stringDiff)
	protest.RequireNext(t, passCode, exitCodes, returnCodeDiff)
	protest.RequireEmpty(t, exitCodes)
	// And that's it
	protest.RequireEmpty(t, calls)
}

func TestRunMutationCatcherFailure(t *testing.T) {
	t.Parallel()

	calls := protest.NewFIFO[string]("calls")
	exitCodes := protest.NewFIFO[returnCodes]("exit codes")

	// Given mutant catcher failure return values from dependencies
	mutantCatcherFails := false
	mutantCatcherFailedCode := returnCodeMutantCatcherFailure
	nilMutationResult := mutationResult{} //nolint:exhaustivestruct,exhaustruct
	theRunner := newMockedRunner(calls, exitCodes, mutantCatcherFails, nilMutationResult)

	// When the func is run
	theRunner.run()

	// Then the program is announced
	protest.RequireNext(t, "announce mutation testing", calls, stringDiff)
	// And the mutant catcher verification is run
	protest.RequireNext(t, "verify mutant catcher passes prior to mutations", calls, stringDiff)
	// And the program exits with 3
	protest.RequireNext(t, "exit", calls, stringDiff)
	protest.RequireNext(t, mutantCatcherFailedCode, exitCodes, returnCodeDiff)
	protest.RequireEmpty(t, exitCodes)
	// And that's it
	protest.RequireEmpty(t, calls)
}

func TestRunUndetectedMutants(t *testing.T) {
	t.Parallel()

	calls := protest.NewFIFO[string]("calls")
	exitCodes := protest.NewFIFO[returnCodes]("exit codes")

	// Given undetected-mutants return values from dependencies
	mutantCatcherPasses := true
	mutantsEscaped := mutationResult{result: experimentResultUndetectedMutants, err: nil}
	failCode := returnCodeFail
	theRunner := newMockedRunner(calls, exitCodes, mutantCatcherPasses, mutantsEscaped)

	// When the func is run
	theRunner.run()

	// Then the program is announced
	protest.RequireNext(t, "announce mutation testing", calls, stringDiff)
	// And the mutant catcher is verified to pass prior to mutations
	protest.RequireNext(t, "verify mutant catcher passes prior to mutations", calls, stringDiff)
	// And the mutations are run
	protest.RequireNext(t, "test mutation types", calls, stringDiff)
	// And the program exits with 1
	protest.RequireNext(t, "exit", calls, stringDiff)
	protest.RequireNext(t, failCode, exitCodes, returnCodeDiff)
	protest.RequireEmpty(t, exitCodes)
	// And that's it
	protest.RequireEmpty(t, calls)
}

func TestRunNoCandidates(t *testing.T) {
	t.Parallel()

	calls := protest.NewFIFO[string]("calls")
	exitCodes := protest.NewFIFO[returnCodes]("exit codes")

	// Given "no candidates" return values from dependencies
	mutantCatcherPasses := true
	mutantsEscaped := mutationResult{result: experimentResultNoCandidatesFound, err: nil}
	failCode := returnCodeNoCandidatesFound
	theRunner := newMockedRunner(calls, exitCodes, mutantCatcherPasses, mutantsEscaped)

	// When the func is run
	theRunner.run()

	// Then the program is announced
	protest.RequireNext(t, "announce mutation testing", calls, stringDiff)
	// And the mutant catcher is verified to pass prior to mutations
	protest.RequireNext(t, "verify mutant catcher passes prior to mutations", calls, stringDiff)
	// And the mutations are run
	protest.RequireNext(t, "test mutation types", calls, stringDiff)
	// And the program exits with no candidates code
	protest.RequireNext(t, "exit", calls, stringDiff)
	protest.RequireNext(t, failCode, exitCodes, returnCodeDiff)
	protest.RequireEmpty(t, exitCodes)
	// And that's it
	protest.RequireEmpty(t, calls)
}

var errMocked = fmt.Errorf("mocked error")

func TestRunError(t *testing.T) {
	t.Parallel()

	calls := protest.NewFIFO[string]("calls")
	exitCodes := protest.NewFIFO[returnCodes]("exit codes")

	// Given "error" return values from dependencies
	mutantCatcherPasses := true
	mutantsEscaped := mutationResult{result: experimentResultError, err: errMocked}
	failCode := returnCodeError
	theRunner := newMockedRunner(calls, exitCodes, mutantCatcherPasses, mutantsEscaped)

	// When the func is run
	theRunner.run()

	// Then the program is announced
	protest.RequireNext(t, "announce mutation testing", calls, stringDiff)
	// And the mutant catcher is verified to pass prior to mutations
	protest.RequireNext(t, "verify mutant catcher passes prior to mutations", calls, stringDiff)
	// And the mutations are run
	protest.RequireNext(t, "test mutation types", calls, stringDiff)
	// And the program exits with error code
	protest.RequireNext(t, "exit", calls, stringDiff)
	protest.RequireNext(t, failCode, exitCodes, returnCodeDiff)
	protest.RequireEmpty(t, exitCodes)
	// And that's it
	protest.RequireEmpty(t, calls)
}

// test out the CLI command
//   announce it
//   announce results
//   if CLI command failed, exit with 3
//   if any error, exit with 2
// search all files under PWD for "true"
//   announce search
//   announce results
//   if none found, exit with 4
//   if any error, exit with 2
// for every matching file search for all instances of "true"
//   announce search
//   announce results
//   if any error, exit with 2
// for every location of "true" run an experiment
//   announce experiment
//   replace "true" with "false" (mutate the candidate)
//   if any error, announce it
//   if any error, exit with 2
//   run the command from the CLI (test whether the command catches the mutant)
//   announce result
//   if any error, exit with 2
//   restore the file to its pre-experiment state
//   if any error, announce it
//   if any error, exit with 2
//   if command from CLI fails to error for the mutant, exit with 1
// announce all mutants found
