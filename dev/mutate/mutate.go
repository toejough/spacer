// Package mutate provides mutation testing functionality.
package main

import "fmt"

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

func main() {
	run(&runDepsMain{
		announceStartingDeps: &announceStartingDepsMain{
			printfunc: func(m string) { fmt.Println(m) },
		},
	})
}

type (
	runDeps interface {
		announceStarting()
		verifyTestsPassWithNoMutants() bool
		testMutations() bool
		announceEnding()
		exit(bool)
	}
	runDepsMain struct {
		announceStartingDeps announceStartingDeps
	}
	announceStartingDeps interface {
		print(string)
	}
	announceStartingDepsMain struct {
		printfunc func(string)
	}
)

func (rdm *runDepsMain) announceStarting() {
	rdm.announceStartingDeps.print("Starting Mutation Testing")
}

func (rdm *runDepsMain) verifyTestsPassWithNoMutants() bool {
	panic("not implemented")
}

func (rdm *runDepsMain) testMutations() bool {
	panic("not implemented")
}

func (rdm *runDepsMain) announceEnding() {
	panic("not implemented")
}

func (rdm *runDepsMain) exit(passes bool) {
	panic("not implemented")
}

func (asdm *announceStartingDepsMain) print(m string) {
	asdm.printfunc(m)
}

// TODO: since methods can't be generic, if a function was supposed to be generic, it could only be included in deps as a function outright, not a method. Which means deps can't be an interface, it _has_ to just be a struct. Which means that for testing, or any other shared state purposes with such a struct, its initialization needs to handle that state via closure.
func run(deps runDeps) {
	deps.announceStarting()
	passes := deps.verifyTestsPassWithNoMutants() && deps.testMutations()
	deps.announceEnding()
	deps.exit(passes)
}
