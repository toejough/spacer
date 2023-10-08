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
	run(&prodRunDeps{})
}

type prodPretestDeps struct{}

// Would need testing for this function if we cared too much about the UI.
func commmonPrintStarting(fname string) func(string) {
	fmt.Printf("%s is starting...\n", fname)

	return func(result string) {
		fmt.Printf("...%s completed with %s\n", fname, result)
	}
}

func (pd *prodPretestDeps) printStarting(fname string) func(string) {
	return commmonPrintStarting(fname)
}
func (pd *prodPretestDeps) fetchPretestCommand() []string { return []string{} }
func (pd *prodPretestDeps) runSubprocess([]string) bool   { return true }

type prodRunDeps struct{}

func (rd *prodRunDeps) printStarting(fname string) func(string) { return commmonPrintStarting(fname) }
func (rd *prodRunDeps) pretest() bool                           { return pretest(&prodPretestDeps{}) }
func (rd *prodRunDeps) testMutations() bool                     { return true }
func (rd *prodRunDeps) exit(code int)                           { os.Exit(code) }
