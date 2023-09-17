// Package mutate provides mutation testing functionality.
package main

import (
	"fmt"
	"os"
)

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

// main runs the program and exits with 0 on success, 1 on failure, 2 on any kind of runtime failure.
func main() {
	fmt.Println("Starting mutation testing")

	pass := run(prodRunDeps())
	if !pass {
		fmt.Println("Mutation testing failed")
		os.Exit(1)
	}

	fmt.Println("Mutation testing passed")
	os.Exit(0)
}

// this function is going to be long... it has all the dependencies.
func prodRunDeps() *runDeps {
	return &runDeps{
		printStarting: func(string) {},
		printDoneWith: func(string) {},
		pretest: func() bool {
			return true
		},
		testMutations: func() bool {
			return true
		},
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
func run(deps *runDeps) bool {
	deps.printStarting("Mutate")
	defer deps.printDoneWith("Mutate")

	return deps.pretest() && deps.testMutations()
}

type (
	runDeps struct {
		printStarting func(what string)
		printDoneWith func(what string)
		pretest       func() bool
		testMutations func() bool
	}
)
