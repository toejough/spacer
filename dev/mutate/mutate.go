// Package mutate provides mutation testing functionality.
package main

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

func main() {
	runner{
		announceMutationTesting:   func() { panic("unimplemented") },
		verifyMutantCatcherPasses: func() mutantCatcherResult { panic("unimplemented") },
		testMutationTypes:         func() mutationResult { panic("unimplemented") },
		announceMutationResults:   func(mutationResult) { panic("unimplemented") },
		exit:                      func(int) { panic("unimplemented") },
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
	testMutationTypesFunc       func() mutationResult
	announceMutationResultsFunc func(mutationResult)
	exitFunc                    func(int)
	runner                      struct {
		announceMutationTesting   announceMutationTestingFunc
		verifyMutantCatcherPasses verifyMutantCatcherPassesFunc
		testMutationTypes         testMutationTypesFunc
		announceMutationResults   announceMutationResultsFunc
		exit                      exitFunc
	}
)

func (r runner) run() {
	r.announceMutationTesting()
	r.verifyMutantCatcherPasses()
	result := r.testMutationTypes()
	r.announceMutationResults(result)
	r.exit(0)
}
