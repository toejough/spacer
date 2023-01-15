package main

import (
	"testing"

	"spacer/dev/protest"

	"pgregory.net/rapid"
)

func TestTestMutationsHappyPath(t *testing.T) {
	t.Parallel()
	// Given inputs/outputs

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
        var testCall testMutationsOnFileCall
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

	return calls, &testMutationsDeps{}
}

type testMutationsDepsMock struct {
	deps  testMutationsDeps
	calls *protest.FIFO[any]
}

type (
	fetchMutationTypesCall protest.CallWithNoArgs[[]mutationType]
	fetchSourceFilesCall   protest.CallWithNoArgs[[]filepath]
    testMutationsOnFileCall protest.Call[testMutationsOnFileArgs, bool]
    testMutationsOnFileArgs struct {
        mutationTypes []mutationType
        path filepath
    }
)
