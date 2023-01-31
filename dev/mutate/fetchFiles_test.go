package main

import (
	"spacer/dev/protest"
	"testing"

	"pgregory.net/rapid"
)

func TestFetchFilesHappyPath(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(test *rapid.T) {
		// Given inputs/outputs
		var result []filepath

		calls, deps := newFetchFilesMock(test)

		// When the function is called
		go func() {
			result = fetchFilesToMutate(deps)

			calls.Close()
		}()

		// Then the paths to include are fetched
        call := calls.MustPop(test)
        protest.MustEqual(test, "fetchPathsToMutate", call.Name)
        protest.MustEqual(test, []any{}, call.PullArgs())

        // When paths are returned
        call.PushReturns(generatedPaths)
        // Then the paths are split into files and directories
        // When the files and directories are returned
        // Then the directories are expanded recursively to the files they contain
        // When the files are returned
        // Then the files are filtered for just .go files
        // When the filtered files are returned
		// Then there are no more calls
		calls.MustConfirmClosed(test)
        // And the files are returned by the fut
		protest.MustEqual(test, files, result)
	})
}


func newFetchFilesMock(test tester) (*protest.FIFO[any], *testMutationsDeps) {
	calls := protest.NewFIFO[any]("calls")

	return calls, &testMutationsDeps{
		fetchMutationTypes: func() []mutationType {
			return protest.ManageCallWithNoArgs[fetchMutationTypesCall](test, calls)
		},
		fetchFilesToMutate: func() []filepath { return protest.ManageCallWithNoArgs[fetchSourceFilesCall](test, calls) },
		testFileMutation: func(f filepath, m []mutationType) bool {
			// return protest.ProxyCall(test, calls, "testFileMutation", f, m).(bool)
			return protest.ManageCall[testFileMutationsCall](test, calls, testFileMutationsArgs{mutationTypes: m, path: f})
		},
	}
}

