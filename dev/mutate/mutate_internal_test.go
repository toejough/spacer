package main

import (
	"spacer/dev/protest"
	"testing"

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

func newMockedDeps(
	calls *protest.FIFO[string],
	exitArgs *protest.FIFO[returnCodes],
	verifyMutantCatcherPassesReturns *protest.FIFO[bool],
	testMutationTypesReturns *protest.FIFO[mutationResult],
) runDeps {
	return runDeps{
		announceMutationTesting: func() { calls.Push("announceMutationTesting") },
		verifyMutantCatcherPasses: func() bool {
			calls.Push("verifyMutantCatcherPasses")
			return verifyMutantCatcherPassesReturns.GetNext()
		},
		testMutationTypes: func() mutationResult {
			calls.Push("testMutationTypes")
			return testMutationTypesReturns.GetNext()
		},
		exit: func(code returnCodes) {
			calls.Push("exit")
			exitArgs.Push(code)
		},
	}
}

func TestRunHappyPath(t *testing.T) {
	t.Parallel()

	// Given Call/Arg/Return FIFOS
	calls := protest.NewFIFO("calls", protest.FIFODeps[string]{
		Differ: stringDiff,
		T:      t,
	})
	exitArgs := protest.NewFIFO("exitArgs", protest.FIFODeps[returnCodes]{Differ: returnCodeDiff, T: t})
	verifyMutantCatcherPassesReturns := protest.NewFIFO("verifyMutantCatcherPassesReturns", protest.FIFODeps[bool]{
		Differ: boolDiff,
		T:      t,
	})
	testMutationTypesReturns := protest.NewFIFO("testMutationTypesReturns", protest.FIFODeps[mutationResult]{
		Differ: mutationResultDiff,
		T:      t,
	})
	deps := newMockedDeps(calls, exitArgs, verifyMutantCatcherPassesReturns, testMutationTypesReturns)

	// When the func is run
	go func() {
		run(deps)
		calls.Close()
	}()

	// Then mutation testing is announced
	calls.RequireNext("announceMutationTesting")
	// And the mutant catcher is tested
	calls.RequireNext("verifyMutantCatcherPasses")

	// When the mutant catcher returns true
	verifyMutantCatcherPassesReturns.Push(true)

	// Then mutation type testing is done
	calls.RequireNext("testMutationTypes")

	// When the testing is all caught
	testMutationTypesReturns.Push(mutationResult{result: experimentResultAllCaught, err: nil})

	// Then the program exits
	calls.RequireNext("exit")
	// and does so with a passing return code
	exitArgs.RequireNext(returnCodePass)
	// and there are no more dependency calls
	calls.RequireClosedAndEmpty()
}

//    t.Parallel()

//    calls := protest.NewFIFO[string]("calls")
//    exitCodes := protest.NewFIFO[returnCodes]("exit codes")

//    // Given mutant catcher failure return values from dependencies
//    mutantCatcherFails := false
//    mutantCatcherFailedCode := returnCodeMutantCatcherFailure
//    nilMutationResult := mutationResult{} //nolint:exhaustivestruct,exhaustruct
//    theRunner := newMockedRunner(calls, exitCodes, mutantCatcherFails, nilMutationResult)

//    // When the func is run
//    theRunner.run()

//    // Then the program is announced
//    protest.RequireNext(t, "announce mutation testing", calls, stringDiff)
//    // And the mutant catcher verification is run
//    protest.RequireNext(t, "verify mutant catcher passes prior to mutations", calls, stringDiff)
//    // And the program exits with 3
//    protest.RequireNext(t, "exit", calls, stringDiff)
//    protest.RequireNext(t, mutantCatcherFailedCode, exitCodes, returnCodeDiff)
//    protest.RequireEmpty(t, exitCodes)
//    // And that's it
//    protest.RequireEmpty(t, calls)
//}

//    t.Parallel()

//    calls := protest.NewFIFO[string]("calls")
//    exitCodes := protest.NewFIFO[returnCodes]("exit codes")

//    // Given undetected-mutants return values from dependencies
//    mutantCatcherPasses := true
//    mutantsEscaped := mutationResult{result: experimentResultUndetectedMutants, err: nil}
//    failCode := returnCodeFail
//    theRunner := newMockedRunner(calls, exitCodes, mutantCatcherPasses, mutantsEscaped)

//    // When the func is run
//    theRunner.run()

//    // Then the program is announced
//    protest.RequireNext(t, "announce mutation testing", calls, stringDiff)
//    // And the mutant catcher is verified to pass prior to mutations
//    protest.RequireNext(t, "verify mutant catcher passes prior to mutations", calls, stringDiff)
//    // And the mutations are run
//    protest.RequireNext(t, "test mutation types", calls, stringDiff)
//    // And the program exits with 1
//    protest.RequireNext(t, "exit", calls, stringDiff)
//    protest.RequireNext(t, failCode, exitCodes, returnCodeDiff)
//    protest.RequireEmpty(t, exitCodes)
//    // And that's it
//    protest.RequireEmpty(t, calls)
//}

//    t.Parallel()

//    calls := protest.NewFIFO[string]("calls")
//    exitCodes := protest.NewFIFO[returnCodes]("exit codes")

//    // Given "no candidates" return values from dependencies
//    mutantCatcherPasses := true
//    mutantsEscaped := mutationResult{result: experimentResultNoCandidatesFound, err: nil}
//    failCode := returnCodeNoCandidatesFound
//    theRunner := newMockedRunner(calls, exitCodes, mutantCatcherPasses, mutantsEscaped)

//    // When the func is run
//    theRunner.run()

//    // Then the program is announced
//    protest.RequireNext(t, "announce mutation testing", calls, stringDiff)
//    // And the mutant catcher is verified to pass prior to mutations
//    protest.RequireNext(t, "verify mutant catcher passes prior to mutations", calls, stringDiff)
//    // And the mutations are run
//    protest.RequireNext(t, "test mutation types", calls, stringDiff)
//    // And the program exits with no candidates code
//    protest.RequireNext(t, "exit", calls, stringDiff)
//    protest.RequireNext(t, failCode, exitCodes, returnCodeDiff)
//    protest.RequireEmpty(t, exitCodes)
//    // And that's it
//    protest.RequireEmpty(t, calls)
//}

// func TestRunError(t *testing.T) {
//    t.Parallel()

//    calls := protest.NewFIFO[string]("calls")
//    exitCodes := protest.NewFIFO[returnCodes]("exit codes")

//    // Given "error" return values from dependencies
//    mutantCatcherPasses := true
//    mutantsEscaped := mutationResult{result: experimentResultError, err: errMocked}
//    failCode := returnCodeError
//    theRunner := newMockedRunner(calls, exitCodes, mutantCatcherPasses, mutantsEscaped)

//    // When the func is run
//    theRunner.run()

//    // Then the program is announced
//    protest.RequireNext(t, "announce mutation testing", calls, stringDiff)
//    // And the mutant catcher is verified to pass prior to mutations
//    protest.RequireNext(t, "verify mutant catcher passes prior to mutations", calls, stringDiff)
//    // And the mutations are run
//    protest.RequireNext(t, "test mutation types", calls, stringDiff)
//    // And the program exits with error code
//    protest.RequireNext(t, "exit", calls, stringDiff)
//    protest.RequireNext(t, failCode, exitCodes, returnCodeDiff)
//    protest.RequireEmpty(t, exitCodes)
//    // And that's it
//    protest.RequireEmpty(t, calls)
//}

//// test out the CLI command
////   announce it
////   announce results
////   if CLI command failed, exit with 3
////   if any error, exit with 2
//// search all files under PWD for "true"
////   announce search
////   announce results
////   if none found, exit with 4
////   if any error, exit with 2
//// for every matching file search for all instances of "true"
////   announce search
////   announce results
////   if any error, exit with 2
//// for every location of "true" run an experiment
////   announce experiment
////   replace "true" with "false" (mutate the candidate)
////   if any error, announce it
////   if any error, exit with 2
////   run the command from the CLI (test whether the command catches the mutant)
////   announce result
////   if any error, exit with 2
////   restore the file to its pre-experiment state
////   if any error, announce it
////   if any error, exit with 2
////   if command from CLI fails to error for the mutant, exit with 1
//// announce all mutants found
