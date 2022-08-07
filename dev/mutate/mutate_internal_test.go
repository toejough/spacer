package main

import (
	"spacer/dev/protest"
	"testing"
)

// announce mutation testing.
func TestWhenProgramStartsAnAnnouncementIsMade(t *testing.T) {
	t.Parallel()

	var call string

	runner{announceMutationTesting: func() { call = "announce mutation testing" }}.run()
	protest.RequireCall(t, "announce mutation testing", call)
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
