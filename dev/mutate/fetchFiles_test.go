package main

import (
	"testing"

	"spacer/dev/protest"

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
		call := calls.MustPopNamed(test, "fetchPathsToMutate")
		protest.MustEqual(test, nil, call.MustPullArgs(test))

		// When paths are returned
		generatedFilesNonGo := rapid.SliceOf(rapid.Custom(drawFilePath)).Draw(test, "generated files nongo")
		generatedFilesGo := rapid.SliceOf(rapid.Custom(drawFilePath)).Draw(test, "generated files go")
		generatedFiles := append(generatedFilesGo, generatedFilesNonGo...)
		generatedDirectories := rapid.SliceOf(rapid.Custom(drawFilePath)).Draw(test, "generated directories")
		generatedPaths := append(generatedFiles, generatedDirectories...)
		call.PushReturns(generatedPaths)
		// Then the paths are split into files and directories
		call = calls.MustPop(test)
		protest.MustEqual(test, "splitFilesAndDirs", call.Name)
		protest.MustEqual(test, generatedPaths, any(call.MustPullArgs(test)).([]filepath))

		// When the files and directories are returned
		call.PushReturns(generatedFiles, generatedDirectories)
		// Then the directories are expanded recursively to the files they contain
		call = calls.MustPop(test)
		protest.MustEqual(test, "recursivelyExpandDirectories", call.Name)
		protest.MustEqual(test, generatedDirectories, any(call.MustPullArgs(test)).([]filepath))

		// When the files are returned
		generatedRecursiveFilesNonGo := rapid.SliceOf(rapid.Custom(drawFilePath)).Draw(test, "generated recursive files nongo")
		generatedRecursiveFilesGo := rapid.SliceOf(rapid.Custom(drawFilePath)).Draw(test, "generated recursive files go")
		generatedRecursiveFiles := append(generatedRecursiveFilesNonGo, generatedRecursiveFilesGo...)
		call.PushReturns(generatedRecursiveFiles)
		// Then the files are filtered for just .go files
		call = calls.MustPop(test)
		protest.MustEqual(test, "onlyGoFiles", call.Name)
		combinedFiles := append(generatedFiles, generatedRecursiveFiles...)
		protest.MustEqual(test, combinedFiles, any(call.MustPullArgs(test)).([]filepath))

		// When the filtered files are returned
		combinedGoFiles := append(generatedFilesGo, generatedRecursiveFilesGo...)
		call.PushReturns(combinedGoFiles, call.MustPullArgs(test))
		// Then there are no more calls
		calls.MustConfirmClosed(test)
		// And the files are returned by the fut
		protest.MustEqual(test, combinedGoFiles, result)
	})
}

func fetchFilesToMutate(deps *fetchFilesDeps) (files []filepath) {
	return deps.fetchPathsToMutate()
}

func newFetchFilesMock(test tester) (*protest.FIFO[protest.AnyCall], *fetchFilesDeps) {
	calls := protest.NewFIFO[protest.AnyCall]("calls")

	return calls, &fetchFilesDeps{
		fetchPathsToMutate: func() []filepath {
			return protest.ProxyCall(test, calls, "fetchPathsToMutate")[0].([]filepath)
		},
	}
}

type fetchFilesDeps struct {
	fetchPathsToMutate func() []filepath
}
