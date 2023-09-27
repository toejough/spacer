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
	run(&prodRunDeps{})
}

type prodPretestDeps struct{}

func (pd *prodPretestDeps) printStarting(string) func(string) { return func(string) {} }
func (pd *prodPretestDeps) fetchPretestCommand() []string     { return []string{} }
func (pd *prodPretestDeps) runSubprocess([]string)            {}

type prodRunDeps struct{}

func (rd *prodRunDeps) printStarting(string) func(string) { return func(string) {} }
func (rd *prodRunDeps) pretest() bool                     { return pretest(&prodPretestDeps{}) }
func (rd *prodRunDeps) testMutations() bool               { return true }
func (rd *prodRunDeps) exit(code int)                     { os.Exit(code) }
