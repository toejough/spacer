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
		verifyMutantCatcherPasses: func() { panic("unimplemented") },
		testMutationTypes:         func() { panic("unimplemented") },
		announceMutationResults:   func() { panic("unimplemented") },
		exit:                      func() { panic("unimplemented") },
	}.run()
}

type (
	announceMutationTestingFunc   func()
	verifyMutantCatcherPassesFunc func()
	testMutationTypesFunc         func()
	announceMutationResultsFunc   func()
	exitFunc                      func()
	runner                        struct {
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
	r.testMutationTypes()
	r.announceMutationResults()
	r.exit()
}
