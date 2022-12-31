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
				print: func(s string) { fmt.Println(s) },
			})
		},
		verifyTestsPassWithNoMutants: func() bool {
			panic("verifyTestsPassWithNoMutants not implemented")
		},
		testMutations: func() bool {
			panic("testMutations not implemented")
		},
		announceEnding: func() {
			panic("announceEnding not implemented")
		},
	}) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

type (
	runDeps struct {
		announceStarting             func()
		verifyTestsPassWithNoMutants func() bool
		testMutations                func() bool
		announceEnding               func()
	}
	announceStartingDeps struct {
		print func(string)
	}
)

func run(deps *runDeps) bool {
	deps.announceStarting()
	passes := deps.verifyTestsPassWithNoMutants() && deps.testMutations()
	deps.announceEnding()

	return passes
}

func announceStarting(deps *announceStartingDeps) {
	deps.print("Starting Mutation Testing")
}
