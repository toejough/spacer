// Package mutate provides mutation testing functionality.
package main

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

// run performs pretesting validation & tests the mutations
//
// Pretesting validation makes sure that before any mutation is run, the test command provided passes.
// If the test command doesn't pass, then that's an error that the user needs to resolve. We owe it to
// them to tell them as much as we can about the test command failure so that they can debug it.
//
// Testing the mutations involves performing mutations one at a time, running the test command, and then
// checking for pass/fail. If the test command passes, that is a failure - it means the mutant was uncaught.
// If the test command fails or errors, we treat that as the mutant being caught.
//
// In either failure case, there's nothing we want to do besides treat it like
// a failure, so the signature of these types is restricted to bools.
func run(deps runDeps) {
	doneFunc := deps.printStarting("Mutate")

	if deps.pretest() && deps.testMutations() {
		doneFunc("Success")
		deps.exit(0)
	} else {
		doneFunc("Failure")
		deps.exit(1)
	}
}

type runDeps interface {
	printStarting(what string) func(string)
	pretest() bool
	testMutations() bool
	exit(int)
}
