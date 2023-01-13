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

func main() {
	if run(&runDeps{
		announceStarting: func() {
			announceStarting(&announceStartingDeps{
				assumePrint: func(s string) { fmt.Println(s) },
			})
		},
		pretest: func() bool {
			return pretest(&pretestDeps{
				announcePretest: func() {
					announcePretest(&announcePretestDeps{
						assumePrint: func(s string) { fmt.Println(s) },
					})
				},
				fetchTestCommand: func() (command, error) {
					panic("fetchTestCommand not implemented")
				},
				runTestCommand: func(command) bool {
					panic("runTestCommand not implemented")
				},
				announcePretestResults: func(bool) {
					panic("announcePretestResults not implemented")
				},
			})
		},
		testMutations: func() bool {
			panic("testMutations not implemented")
		},
		announceEnding: func(bool) {
			panic("announceEnding not implemented")
		},
	}) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func announcePretest(announcePretestDeps *announcePretestDeps) {
	announcePretestDeps.assumePrint("Starting Pretesting")
}

type (
	runDeps struct {
		announceStarting func()
		pretest          func() bool
		testMutations    func() bool
		announceEnding   func(bool)
	}
	announceStartingDeps struct {
		assumePrint func(string)
	}
	announcePretestDeps struct {
		assumePrint func(string)
	}
	command     string
	pretestDeps struct {
		announcePretest        func()
		fetchTestCommand       func() (command, error)
		runTestCommand         func(command) bool
		announcePretestResults func(bool)
	}
)

func run(deps *runDeps) bool {
	deps.announceStarting()
	passes := deps.pretest() && deps.testMutations()
	deps.announceEnding(passes)

	return passes
}

func announceStarting(deps *announceStartingDeps) {
	deps.assumePrint("Starting Mutation Testing")
}

func pretest(deps *pretestDeps) bool {
	deps.announcePretest()

	c, err := deps.fetchTestCommand()
	if err != nil {
		return false
	}

	result := deps.runTestCommand(c)
	deps.announcePretestResults(result)

	return result
}
