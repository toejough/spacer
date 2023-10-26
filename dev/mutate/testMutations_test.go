package main

import (
	"testing"
	"time"
)

func TestTestMutationsHappyPath(t *testing.T) {
	// Given Test Needs
	tester := protest.NewTester(t)

	// Given Inputs, dependencies, and data from dependencies
	deps := testMutationsDepsForTest{tester: tester}
	goFiles := []string{}
	mutators := []string{}
	testCommand := []string{}
	// Given Outputs
	result := true

	// When Func is run
	tester.Start(testMutations)

	// Then find the go files, & alert the user
	tester.AssertCalledNext(deps.FetchGoFiles).InjectReturns(goFiles)
	tester.AssertCalledNext(deps.CommunicateGoFiles, goFiles)
	// Then find the mutators, & alert the user
	tester.AssertCalledNext(deps.FetchMutators).InjectReturns(mutators)
	tester.AssertCalledNext(deps.CommunicateMutators, mutators)
	// Then get the test command
	tester.AssertCalledNext(deps.FetchTestCommand).InjectReturns(testCommand)
	// Then for each file & mutator, find the candidates & alert the user
	tester.AssertCalledNext(deps.IdentifyCandidates, goFiles, mutators).InjectReturns(candidates)
	// Then for each candidate, cache the file, mutate, alert the user about the difference & run the test command
	tester.AssertCalledNext(deps.TestCandidates, candidates).InjectReturns(candidateResults)
	// Then alert the user about the results of the test
	tester.AssertCalledNext(deps.CommunicateCandidateResults, candidateResults)

	// Then assert func is done
	tester.AssertDoneWithin(time.Second)

	// Then if any escapees, fail (otherwise pass)
	tester.AssertReturned(result)
}
