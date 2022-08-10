package main

import (
	"spacer/dev/protest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func stringDiff(e, a string) string {
	return cmp.Diff(e, a)
}

// announce mutation testing.
func TestWhenProgramStartsAnAnnouncementIsMade(t *testing.T) {
	t.Parallel()

	calls := protest.NewFIFO[string]("calls")

	runner{
		announceMutationTesting:   func() { calls.Push("announce mutation testing") },
		verifyMutantCatcherPasses: func() { calls.Push("verify mutant catcher passes prior to mutations") },
		testMutationTypes:         func() { calls.Push("test mutation types") },
		announceMutationResults:   func() { calls.Push("announce mutation results") },
		exit:                      func() { calls.Push("exit") },
	}.run()

	protest.RequireNext(t, "announce mutation testing", calls, stringDiff)
	protest.RequireNext(t, "verify mutant catcher passes prior to mutations", calls, stringDiff)
	protest.RequireNext(t, "test mutation types", calls, stringDiff)
	protest.RequireNext(t, "announce mutation results", calls, stringDiff)
	protest.RequireNext(t, "exit", calls, stringDiff)
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
