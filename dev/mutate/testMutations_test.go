package main

import (
	"testing"
	"time"
)

func TestTestMutationsHappyPath(t *testing.T) {
	// Given Test Needs
	tester := protest.NewTester(t)
	// Given dependencies
	deps := testMutationsDepsForTest{tester: tester}
	// Given data from dependencies
	goFiles := []string{}
	mutators := []string{}
	testCommand := []string{}
	candidates := []string{}
	candidateResults := []string{}
	// Given Outputs
	result := true

	// When Func is run
	tester.Start(testMutations)

	// Then find the go files
	tester.AssertCalledNext(deps.FetchGoFiles).InjectReturns(goFiles)
	// Then find the mutators
	tester.AssertCalledNext(deps.FetchMutators).InjectReturns(mutators)
	// Then get the test command
	tester.AssertCalledNext(deps.FetchTestCommand).InjectReturns(testCommand)
	// Then for each file & mutator, find the candidates & alert the user
	tester.AssertCalledNext(deps.IdentifyCandidates, goFiles, mutators).InjectReturns(candidates)
	// Then test each candidate
	tester.AssertCalledNext(deps.TestCandidates, candidates).InjectReturns(candidateResults)
	// Then assert func is done
	tester.AssertDoneWithin(time.Second)
	// Then if any escapees, fail (otherwise pass)
	tester.AssertReturned(result)
}
