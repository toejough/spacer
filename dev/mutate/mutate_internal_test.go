package main

import (
	"spacer/dev/protest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func stringDiff(e, a string) string {
	return cmp.Diff(e, a)
}

func mutationResultDiff(e, a mutationResult) string {
	// struct passed here for type, not for data
	return cmp.Diff(e, a, cmp.AllowUnexported(mutationResult{})) //nolint:exhaustivestruct,exhaustruct
}

func intDiff(e, a int) string {
	return cmp.Diff(e, a)
}

// announce mutation testing.
func TestRunHappyPath(t *testing.T) {
	t.Parallel()

	calls := protest.NewFIFO[string]("calls")
	mutantCatcherPasses := mutantCatcherResult{pass: true, err: nil}
	mutationResults := protest.NewFIFO[mutationResult]("mutation results")
	allMutantsCaught := mutationResult{allCaught: true, err: nil}
	exitCodes := protest.NewFIFO[int]("exit codes")
	passCode := 0

	runner{
		announceMutationTesting: func() { calls.Push("announce mutation testing") },
		verifyMutantCatcherPasses: func() mutantCatcherResult {
			calls.Push("verify mutant catcher passes prior to mutations")
			return mutantCatcherPasses
		},
		testMutationTypes: func() mutationResult {
			calls.Push("test mutation types")
			return allMutantsCaught
		},
		announceMutationResults: func(r mutationResult) {
			calls.Push("announce mutation results")
			mutationResults.Push(r)
		},
		exit: func(code int) {
			calls.Push("exit")
			exitCodes.Push(code)
		},
	}.run()

	protest.RequireNext(t, "announce mutation testing", calls, stringDiff)
	protest.RequireNext(t, "verify mutant catcher passes prior to mutations", calls, stringDiff)
	protest.RequireNext(t, "test mutation types", calls, stringDiff)

	protest.RequireNext(t, "announce mutation results", calls, stringDiff)
	protest.RequireNext(t, allMutantsCaught, mutationResults, mutationResultDiff)
	protest.RequireEmpty(t, mutationResults)

	protest.RequireNext(t, "exit", calls, stringDiff)
	protest.RequireNext(t, passCode, exitCodes, intDiff)
	protest.RequireEmpty(t, exitCodes)

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
