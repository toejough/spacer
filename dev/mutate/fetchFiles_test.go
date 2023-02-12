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
		call := calls.MustPopNamed(test, "fetchPathsToMutate")
		protest.MustEqual(test, nil, call.MustPullArgs(test))

		// When paths are returned
		generatedPaths := rapid.SliceOf(rapid.Custom(drawFilePath)).Draw(test, "generated paths")
		call.PushReturns(generatedPaths)
		// Then the paths are split into files and directories
		call = calls.MustPopNamed(test, "splitFilesAndDirs")
		protest.MustEqual(test, []any{generatedPaths}, call.MustPullArgs(test))

		// When the files and directories are returned
		generatedFiles := rapid.SliceOf(rapid.Custom(drawFilePath)).Draw(test, "generated files")
		generatedDirectories := rapid.SliceOf(rapid.Custom(drawFilePath)).Draw(test, "generated directories")
		call.PushReturns(generatedFiles, generatedDirectories)
		// Then the directories are expanded recursively to the files they contain
		call = calls.MustPopNamed(test, "recursivelyExpandDirectories")
		protest.MustEqual(test, []any{generatedDirectories}, call.MustPullArgs(test))

		// When the files are returned
		generatedRecursiveFiles := rapid.SliceOf(rapid.Custom(drawFilePath)).Draw(test, "generated recursive files")
		call.PushReturns(generatedRecursiveFiles)
		// Then the files are filtered for just .go files
		call = calls.MustPopNamed(test, "filterToGoFiles")
		combinedFiles := combine(generatedFiles, generatedRecursiveFiles)
		protest.MustEqual(test, []any{combinedFiles}, call.MustPullArgs(test))

		// When the filtered files are returned
		combinedFilesGo := rapid.SliceOf(rapid.Custom(drawFilePath)).Draw(test, "combined files go")
		call.PushReturns(combinedFilesGo)
		// Then there are no more calls
		calls.MustConfirmClosed(test)
		// And the files are returned by the fut
		protest.MustEqual(test, combinedFilesGo, result)
	})
}

func combine(a, b []filepath) []filepath {
	combined := make([]filepath, len(a)+len(b))
	copy(combined, a)

	for i, item := range b {
		combined[len(a)+i] = item
	}

	return combined
}

func fetchFilesToMutate(deps *fetchFilesDeps) (filesToMutate []filepath) {
	paths := deps.fetchPathsToMutate()
	files, dirs := deps.splitFilesAndDirs(paths)
	expandedFiles := deps.recursivelyExpandDirectories(dirs)
	allFiles := combine(files, expandedFiles)

	return deps.filterToGoFiles(allFiles)
}

func newFetchFilesMock(test tester) (*protest.FIFO[protest.AnyCall], *fetchFilesDeps) {
	calls := protest.NewFIFO[protest.AnyCall]("calls")

	return calls, &fetchFilesDeps{
		fetchPathsToMutate: func() []filepath {
			returns := protest.ProxyCall(test, calls, "fetchPathsToMutate")
			return protest.MustUnwrapTo1[[]filepath](test, returns)
		},
		splitFilesAndDirs: func(paths []filepath) (files, dirs []filepath) {
			returns := protest.ProxyCall(test, calls, "splitFilesAndDirs", paths)
			return protest.MustUnwrapTo2[[]filepath, []filepath](test, returns)
		},
		recursivelyExpandDirectories: func(dirs []filepath) (files []filepath) {
			returns := protest.ProxyCall(test, calls, "recursivelyExpandDirectories", dirs)
			return protest.MustUnwrapTo1[[]filepath](test, returns)
		},
		filterToGoFiles: func(files []filepath) (goFiles []filepath) {
			returns := protest.ProxyCall(test, calls, "filterToGoFiles", files)
			return protest.MustUnwrapTo1[[]filepath](test, returns)
		},
	}
}

type fetchFilesDeps struct {
	fetchPathsToMutate           func() []filepath
	splitFilesAndDirs            func(paths []filepath) (files, dirs []filepath)
	recursivelyExpandDirectories func(dirs []filepath) (files []filepath)
	filterToGoFiles              func(files []filepath) (goFiles []filepath)
}
