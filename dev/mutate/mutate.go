// Package mutate provides mutation testing functionality.
package main

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

func main() {
	run(&runDepsMain{})
}

type (
	runDeps interface {
		announceStarting()
		verifyTestsPassWithNoMutants() bool
		testMutations() bool
		announceEnding()
		exit(bool)
	}
	runDepsMain struct{}
)

func (rdm *runDepsMain) announceStarting() {
	panic("not implemented")
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

func run(deps runDeps) {
	deps.announceStarting()

	passes := deps.verifyTestsPassWithNoMutants()
	if !passes {
		deps.announceEnding()
		deps.exit(false)

		return
	}

	passes = deps.testMutations()
	deps.announceEnding()
	deps.exit(passes)
}
