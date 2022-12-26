// Package mutate provides mutation testing functionality.
package main

import "fmt"

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

func main() {
	run(runDeps{
		announceMutationTesting:   func() { fmt.Println("Starting mutation testing") },
		verifyMutantCatcherPasses: nil,
		testMutationTypes:         nil,
		exit:                      nil,
	})
}

type (
	announceMutationTestingFunc   func()
	verifyMutantCatcherPassesFunc func() bool
	experimentResult              int
	mutationResult                struct {
		result experimentResult
		err    error
	}
	testMutationTypesFunc func() mutationResult
	returnCodes           int
	exitFunc              func(returnCodes)
	runDeps               struct {
		announceMutationTesting   announceMutationTestingFunc
		verifyMutantCatcherPasses verifyMutantCatcherPassesFunc
		testMutationTypes         testMutationTypesFunc
		exit                      exitFunc
	}
)

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
