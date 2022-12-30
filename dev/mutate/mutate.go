// Package mutate provides mutation testing functionality.
package main

import "fmt"

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

func main() {
	run(&runDeps{
		announceStarting: func() {
			announceStarting(announceStartingDeps{
				print: func(s string) { fmt.Println(s) },
			})
		},
		verifyTestsPassWithNoMutants: func() bool {
			panic("not implemented")
		},
		testMutations: func() bool {
			panic("not implemented")
		},
		announceEnding: func() {
			panic("not implemented")
		},
		exit: func(passes bool) {
			panic("not implemented")
		},
	})
}

type (
	runDeps struct {
		announceStarting             func()
		verifyTestsPassWithNoMutants func() bool
		testMutations                func() bool
		announceEnding               func()
		exit                         func(bool)
	}
	announceStartingDeps struct {
		print func(string)
	}
)

func run(deps *runDeps) {
	deps.announceStarting()
	passes := deps.verifyTestsPassWithNoMutants() && deps.testMutations()
	deps.announceEnding()
	deps.exit(passes)
}

func announceStarting(deps announceStartingDeps) {
	deps.print("Starting Mutation Testing")
}
