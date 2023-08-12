package main

import (
	"spacer/dev/protest"
	"testing"

	"pgregory.net/rapid"
)

// TODO - refactor all old tests to use the new protest funcs
// TODO - kill off the unused protest funcs
// TODO - tidy vertical call flows across all files.
// TODO - instead of MustPopNamed on all Fifos that has a type assertion, make a special calls FIFO with a MustPop.
// TODO - actually make one-shots their own type.
// TODO - test getting mutation types
// name, AST predicate func, AST modification func. Check out https://eli.thegreenplace.net/2021/rewriting-go-source-code-with-ast-tooling/
// TODO - test running mutation types.
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

func newFetchFilesMock(test tester) (*protest.FIFO[protest.AnyCall], *fetchFilesDeps) {
	calls := protest.NewFIFO[protest.AnyCall]("calls")

	return calls, &fetchFilesDeps{
		fetchPathsToMutate: func() []filepath {
			return protest.ProxyCallR1[[]filepath](test, calls, "fetchPathsToMutate")
		},
		splitFilesAndDirs: func(paths []filepath) (files, dirs []filepath) {
			return protest.ProxyCallR2[[]filepath, []filepath](test, calls, "splitFilesAndDirs", paths)
		},
		recursivelyExpandDirectories: func(dirs []filepath) (files []filepath) {
			return protest.ProxyCallR1[[]filepath](test, calls, "recursivelyExpandDirectories", dirs)
		},
		filterToGoFiles: func(files []filepath) (goFiles []filepath) {
			return protest.ProxyCallR1[[]filepath](test, calls, "filterToGoFiles", files)
		},
	}
}
