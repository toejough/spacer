// Package mutate provides mutation testing functionality.
package main

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

func main() {
	runner{
		announceMutationTesting:   nil,
		verifyMutantCatcherPasses: nil,
		testMutationTypes:         nil,
		exit:                      nil,
	}.run()
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
	runner                struct {
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
	returnCodeFail
	returnCodeMutantCatcherFailure
	returnCodeNoCandidatesFound
)

func (r runner) run() {
	r.announceMutationTesting()

	passes := r.verifyMutantCatcherPasses()
	if !passes {
		r.exit(returnCodeMutantCatcherFailure)
		return
	}

	results := r.testMutationTypes()
	switch results.result {
	case experimentResultAllCaught:
		r.exit(returnCodePass)
		return
	case experimentResultUndetectedMutants:
		r.exit(returnCodeFail)
		return
	case experimentResultNoCandidatesFound:
		r.exit(returnCodeNoCandidatesFound)
		return
	case experimentResultError:
		r.exit(returnCodeError)
		return
	}
}
