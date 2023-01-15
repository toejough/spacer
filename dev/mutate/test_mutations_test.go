package main

import (
	"testing"

	"spacer/dev/protest"
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

	// Then the mutation inputs are fetched:
	// * the mutation types
	// * the source files to mutate
	first := calls.MustPop(t)
	second := calls.MustPop(t)
	var mutationTypesCall fetchMutationTypesCall
	var sourceFilesCall fetchSourceFilesCall
	var ok bool

	// TODO this is the most naive way to do this - make a function that:
	// * takes a slice of any and a target
	// * sets the target if it can from something in the slice
	// * returns a slice of the remaining items and an OK indicating whether or not the target was set.
	mutationTypesCall, ok = first.(fetchMutationTypesCall)
	if !ok {
		mutationTypesCall, ok = second.(fetchMutationTypesCall)
		if !ok {
			t.Fatalf("neither of the popped calls were for fetching the mutation types: %v, %v", first, second)
		}
	}

	sourceFilesCall, ok = first.(fetchSourceFilesCall)
	if !ok {
		sourceFilesCall, ok = second.(fetchSourceFilesCall)
		if !ok {
			t.Fatalf("neither of the popped calls were for fetching the source files: %v, %v", first, second)
		}
	}

	// When the inputs have been fetched
	mutationTypes := []mutationType{} // TODO rapid test this
	mutationTypesCall.ReturnOneShot.Push(mutationTypes)
	sourceFiles := []filepath{} // TODO rapid test this
	sourceFilesCall.ReturnOneShot.Push(sourceFiles)

	// Then each file is tested for all mutation types
	for i := 0; i < len(sourceFiles); i++ {
		var testCall testFileMutationsCall
		calls.MustPopAs(t, testCall)
		protest.MustEqual(t, testCall.Args.mutationTypes, mutationTypes)
		if !contains(sourceFiles, testCall.Args.path) {
			t.Fatalf("no call expected for the given path: %s", testCall.Args.path)
		}
		sourceFiles = remove(sourceFiles, testCall.Args.path)

		// When all tests pass
		testCall.ReturnOneShot.Push(true)
	}

	// Then passing status is returned
	protest.MustEqual(t, true, result)
	// and there are no more calls
	calls.MustConfirmClosed(t)
}

func newTestMutationsMock(test tester) (*protest.FIFO[any], *testMutationsDeps) {
	calls := protest.NewFIFO[any]("calls")

	return calls, &testMutationsDeps{
        fetchMutationTypes: func() []mutationType { return protest.ManageCallWithNoArgs[fetchMutationTypesCall](test, calls)},
        fetchFilesToMutate: func() []filepath {return protest.ManageCallWithNoArgs[fetchSourceFilesCall](test, calls)},
        testFileMutation: func(f filepath, m []mutationType) bool {return protest.ManageCall[testFileMutationsCall](test, calls, testFileMutationsArgs{mutationTypes: m, path: f})},
    }
}

func contains[I any](slice []I, item I) bool {
	for _, i := range slice {
		if protest.Equal(item, i) != nil {
			return true
		}
	}

	return false
}

func remove[I any](slice []I, item I) (newSlice []I) {
	index := 0
	for index = range slice {
		if protest.Equal(item, slice[index]) == nil {
			if index+1 < len(newSlice) {
				newSlice = append(newSlice, slice[index+1:]...)
			}
			break
		}
		newSlice = append(newSlice, slice[index])
	}

	return newSlice
}

type (
	testMutationsDepsMock struct {
		deps  testMutationsDeps
		calls *protest.FIFO[any]
	}
	fetchMutationTypesCall protest.CallWithNoArgs[[]mutationType]
	fetchSourceFilesCall   protest.CallWithNoArgs[[]filepath]
	testFileMutationsCall  protest.Call[testFileMutationsArgs, bool]
	testFileMutationsArgs  struct {
		mutationTypes []mutationType
		path          filepath
	}
)
