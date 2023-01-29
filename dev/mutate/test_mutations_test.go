package main

import (
	"spacer/dev/protest"
	"testing"

	"pgregory.net/rapid"
)

func TestTestMutationsHappyPath(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(test *rapid.T) {
		// Given inputs/outputs
		var result bool

		calls, deps := newTestMutationsMock(test)

		// When the function is called
		go func() {
			result = testMutations(deps)

			calls.Close()
		}()

		// Then the mutation types are fetched
		var mutationTypesCall fetchMutationTypesCall

		calls.MustPopAs(test, &mutationTypesCall)

		// When the mutation types are returned
		mutationTypes := rapid.SliceOfN(rapid.Make[mutationType](), 1, -1).Draw(test, "mutationTypes")
		mutationTypesCall.ReturnOneShot.Push(mutationTypes)

		// Then the source file paths are fetched
		var sourceFilesCall fetchSourceFilesCall

		calls.MustPopAs(test, &sourceFilesCall)

		// When the source file paths are returned
		sourceFiles := rapid.SliceOfN(rapid.Custom(drawFilePath), 1, -1).Draw(test, "filepaths")
		sourceFilesCall.ReturnOneShot.Push(sourceFiles)

		// Then each file is tested for all mutation types
		for _, fp := range sourceFiles {
			var testCall testFileMutationsCall

			calls.MustPopAs(test, &testCall)
			protest.MustEqual(test, testCall.Args, testFileMutationsArgs{mutationTypes: mutationTypes, path: fp})

			// When all tests pass
			testCall.ReturnOneShot.Push(true)
		}

		// Then there are no more calls
		calls.MustConfirmClosed(test)
		// and a passing status is returned
		protest.MustEqual(test, true, result)
	})
}

func drawFilePath(t *rapid.T) filepath {
	return filepath(rapid.String().Draw(t, "filepath"))
}

func TestTestMutationsNoFiles(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(test *rapid.T) {
		// Given inputs/outputs
		var result bool

		calls, deps := newTestMutationsMock(test)

		// When the function is called
		go func() {
			result = testMutations(deps)

			calls.Close()
		}()

		// Then the mutation types are fetched
		var mutationTypesCall fetchMutationTypesCall

		calls.MustPopAs(test, &mutationTypesCall)

		// When the mutation types are returned
		mutationTypes := rapid.SliceOfN(rapid.Make[mutationType](), 1, -1).Draw(test, "mutationTypes")
		mutationTypesCall.ReturnOneShot.Push(mutationTypes)

		// Then the source file paths are fetched
		var sourceFilesCall fetchSourceFilesCall

		calls.MustPopAs(test, &sourceFilesCall)

		// When no source file paths are returned
		sourceFiles := []filepath{}
		sourceFilesCall.ReturnOneShot.Push(sourceFiles)

		// Then there are no more calls
		calls.MustConfirmClosed(test)
		// and a failing status is returned
		protest.MustEqual(test, false, result)
	})
}

func TestTestMutationsNoMutationTypes(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(test *rapid.T) {
		// Given inputs/outputs
		var result bool

		calls, deps := newTestMutationsMock(test)

		// When the function is called
		go func() {
			result = testMutations(deps)

			calls.Close()
		}()

		// Then the mutation types are fetched
		var mutationTypesCall fetchMutationTypesCall

		calls.MustPopAs(test, &mutationTypesCall)

		// When no mutation types are returned
		mutationTypes := []mutationType{}
		mutationTypesCall.ReturnOneShot.Push(mutationTypes)

		// Then there are no more calls
		calls.MustConfirmClosed(test)
		// and a failing status is returned
		protest.MustEqual(test, false, result)
	})
}

func TestTestMutationsUncaught(t *testing.T) {
	t.Parallel()

	rapid.Check(t, func(test *rapid.T) {
		// Given inputs/outputs
		var result bool

		calls, deps := newTestMutationsMock(test)

		// When the function is called
		go func() {
			result = testMutations(deps)

			calls.Close()
		}()

		// Then the mutation types are fetched
		var mutationTypesCall fetchMutationTypesCall

		calls.MustPopAs(test, &mutationTypesCall)

		// When the mutation types are returned
		mutationTypes := rapid.SliceOfN(rapid.Make[mutationType](), 1, -1).Draw(test, "mutationTypes")
		mutationTypesCall.ReturnOneShot.Push(mutationTypes)

		// Then the source file paths are fetched
		var sourceFilesCall fetchSourceFilesCall

		calls.MustPopAs(test, &sourceFilesCall)

		// When the source file paths are returned
		sourceFiles := rapid.SliceOfN(rapid.Custom(drawFilePath), 1, -1).Draw(test, "filepaths")
		sourceFilesCall.ReturnOneShot.Push(sourceFiles)
		numSourceFiles := len(sourceFiles)
		testResults := rapid.SliceOfN(
			rapid.SampledFrom([]bool{true, false}), numSourceFiles, numSourceFiles,
		).Filter(atLeastOneFalse).Draw(test, "mutation test result")

		// Then each file is tested for all mutation types
		for index, fp := range sourceFiles {
			var testCall testFileMutationsCall

			calls.MustPopAs(test, &testCall)
			protest.MustEqual(test, testCall.Args, testFileMutationsArgs{mutationTypes: mutationTypes, path: fp})

			// When all tests pass
			testCall.ReturnOneShot.Push(testResults[index])

			if !testResults[index] {
				break
			}
		}

		// Then there are no more calls
		calls.MustConfirmClosed(test)
		// and a passing status is returned
		protest.MustEqual(test, false, result)
	})
}

func atLeastOneFalse(items []bool) bool {
	for _, i := range items {
		if !i {
			return true
		}
	}

	return false
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
