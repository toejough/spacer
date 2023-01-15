package main

import (
	"spacer/dev/protest"
	"testing"
)

func TestTestMutationsHappyPath(t *testing.T) {
	t.Parallel()

	// Given inputs/outputs
	var result bool

	calls, deps := newTestMutationsMock(t)

	// When the function is called
	go func() {
		result = testMutations(deps)

		calls.Close()
	}()

	// Then the mutation types are fetched
	var mutationTypesCall fetchMutationTypesCall

	calls.MustPopAs(t, &mutationTypesCall)

	// When the mutation types are returned
	mutationTypes := []mutationType{} // TODO rapid test this
	mutationTypesCall.ReturnOneShot.Push(mutationTypes)

	// Then the source file paths are fetched
	var sourceFilesCall fetchSourceFilesCall

	calls.MustPopAs(t, &sourceFilesCall)

	// When the source file paths are returned
	sourceFiles := []filepath{} // TODO rapid test this
	sourceFilesCall.ReturnOneShot.Push(sourceFiles)

	// Then each file is tested for all mutation types
	for _, fp := range sourceFiles {
		var testCall testFileMutationsCall

		calls.MustPopAs(t, testCall)
		protest.MustEqual(t, testCall.Args, testFileMutationsArgs{mutationTypes: mutationTypes, path: fp})

		// When all tests pass
		testCall.ReturnOneShot.Push(true)
	}

	// Then there are no more calls
	calls.MustConfirmClosed(t)
	// and a passing status is returned
	protest.MustEqual(t, true, result)
}

func newTestMutationsMock(test tester) (*protest.FIFO[any], *testMutationsDeps) {
	calls := protest.NewFIFO[any]("calls")

	return calls, &testMutationsDeps{
		fetchMutationTypes: func() []mutationType {
			return protest.ManageCallWithNoArgs[fetchMutationTypesCall](test, calls)
		},
		fetchFilesToMutate: func() []filepath { return protest.ManageCallWithNoArgs[fetchSourceFilesCall](test, calls) },
		testFileMutation: func(f filepath, m []mutationType) bool {
			return protest.ManageCall[testFileMutationsCall](test, calls, testFileMutationsArgs{mutationTypes: m, path: f})
		},
	}
}

type (
	fetchMutationTypesCall protest.CallWithNoArgs[[]mutationType]
	fetchSourceFilesCall   protest.CallWithNoArgs[[]filepath]
	testFileMutationsCall  protest.Call[testFileMutationsArgs, bool]
	testFileMutationsArgs  struct {
		mutationTypes []mutationType
		path          filepath
	}
)
