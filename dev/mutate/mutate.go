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
	verifyMutantCatcherPassesFunc func() mutantCatcherResult
	mutationResult                struct {
		allCaught bool
		err       error
	}
	mutantCatcherResult struct {
		pass bool
		err  error
	}
	testMutationTypesFunc func() mutationResult
	exitFunc              func(returnCodes)
	runner                struct {
		announceMutationTesting   announceMutationTestingFunc
		verifyMutantCatcherPasses verifyMutantCatcherPassesFunc
		testMutationTypes         testMutationTypesFunc
		exit                      exitFunc
	}
)

type returnCodes int

const (
	returnCodePass returnCodes = iota
	returnCodeMutantCatcherFailure
)

func (r runner) run() {
	r.announceMutationTesting()

	results := r.verifyMutantCatcherPasses()
	if !results.pass {
		r.exit(returnCodeMutantCatcherFailure)
		return
	}

	r.testMutationTypes()
	r.exit(returnCodePass)
}
