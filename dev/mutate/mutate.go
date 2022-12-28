// Package mutate provides mutation testing functionality.
package main

import "fmt"

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

func main() {
	run(&runDepsMain{})
}

type (
	experimentResult int
	mutationResult   struct {
		result experimentResult
		err    error
	}
	returnCodes int
	runDeps     interface {
		announceMutationTesting()
		verifyMutantCatcherPasses() bool
		testMutationTypes() mutationResult
		exit(returnCodes)
	}
	runDepsMain struct{}
)

func (rdm *runDepsMain) announceMutationTesting() {
	fmt.Println("Starting mutation testing")
}

func (rdm *runDepsMain) verifyMutantCatcherPasses() bool {
	panic("not implemented")
}

func (rdm *runDepsMain) testMutationTypes() mutationResult {
	panic("not implemented")
}

func (rdm *runDepsMain) exit(rc returnCodes) {
	panic("not implemented")
}

const (
	experimentResultAllCaught experimentResult = iota
	experimentResultUndetectedMutants
	experimentResultNoCandidatesFound
	experimentResultError
)

const (
	returnCodePass returnCodes = iota
	returnCodeFail
	returnCodeError
	returnCodeMutantCatcherFailure
	returnCodeNoCandidatesFound
)

func run(deps runDeps) {
	deps.announceMutationTesting()

	passes := deps.verifyMutantCatcherPasses()
	if !passes {
		deps.exit(returnCodeMutantCatcherFailure)
		return
	}

	results := deps.testMutationTypes()
	switch results.result {
	case experimentResultAllCaught:
		deps.exit(returnCodePass)
		return
	case experimentResultUndetectedMutants:
		deps.exit(returnCodeFail)
		return
	case experimentResultNoCandidatesFound:
		deps.exit(returnCodeNoCandidatesFound)
		return
	case experimentResultError:
		deps.exit(returnCodeError)
		return
	}
}
