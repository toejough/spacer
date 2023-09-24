// Package mutate provides mutation testing functionality.
package main

import "os"

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

// main runs the program and exits with 0 on success, 1 on failure, 2 on any kind of runtime failure.
func main() {
	run(prodRunDeps())
}

type prodPretestDeps struct{}

func (pd *prodPretestDeps) printStarting(string) func(string) {
	return func(string) {}
}

func (pd *prodPretestDeps) fetchPretestCommand() []string {
	return []string{}
}

func (pd *prodPretestDeps) runSubprocess([]string) {}

// this function is going to be long... it has all the dependencies.
func prodRunDeps() *runDeps {
	return &runDeps{
		printStarting: func(string) func(string) { return func(string) {} },
		pretest: func() bool {
			return pretest(&prodPretestDeps{})
		},
		testMutations: func() bool {
			return true
		},
		exit: func(code int) { os.Exit(code) },
	}
}

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
func run(deps *runDeps) {
	doneFunc := deps.printStarting("Mutate")

	if deps.pretest() && deps.testMutations() {
		doneFunc("Success")
		deps.exit(0)
	} else {
		doneFunc("Failure")
		deps.exit(1)
	}
}

type (
	runDeps struct {
		printStarting func(what string) func(string)
		pretest       func() bool
		testMutations func() bool
		exit          func(int)
	}
)
