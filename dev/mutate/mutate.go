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
	experimentResult int
	mutationResult   struct {
		result experimentResult
		err    error
	}
	returnCodes int
	runDeps     interface {
		verifyTestsPassWithNoMutants() error
		testMutationTypes() mutationResult
		exit(returnCodes)
	}
	runDepsMain struct{}
)

func (rdm *runDepsMain) verifyTestsPassWithNoMutants() error {
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
	returnCodeTestsFailWithNoMutations
	returnCodeNoCandidatesFound
)

func run(deps runDeps) {
	err := deps.verifyTestsPassWithNoMutants()
	if err != nil {
		deps.exit(returnCodeTestsFailWithNoMutations)
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
