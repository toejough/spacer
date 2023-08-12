// Package mutate provides mutation testing functionality.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

func main() {
	fmt.Println("Starting mutation testing")

	if run(prodRunDeps()) {
		fmt.Println("Mutation testing passed")
		os.Exit(0)
	} else {
		fmt.Println("Mutation testing failed")
		os.Exit(1)
	}
}

// this function is going to be long... it has all the dependencies.
func prodRunDeps() *runDeps { //nolint:funlen
	return &runDeps{
		pretest: func() bool {
			fmt.Println("Starting pretesting")
			results := pretest(&pretestDeps{
				fetchTestCommand: func() command {
					fmt.Println("Fetching test command")
					if len(os.Args) < 2 { //nolint:gomnd
						fmt.Println("no test command provided on CLI")
						return ""
					}
					c := os.Args[1]
					fmt.Printf("Fetched '%s' as the command\n", c)

					return command(c)
				},
				runTestCommand: func(comm command) bool {
					fmt.Println("Running test command")
					parts := strings.Split(string(comm), " ")
					commObj := exec.Command(parts[0], parts[1:]...) //nolint:gosec
					output, err := commObj.Output()
					if err != nil {
						fmt.Printf("Test command failed: %v\n", err)
						return false
					}

					fmt.Printf("Test command passed: %s\n", output)

					return true
				},
			})
			fmt.Printf("Pretest passed? %t\n", results)

			return results
		},
		testMutations: func() bool {
			fmt.Println("Testing mutations")
			if !testMutations(&testMutationsDeps{
				fetchMutationTypes: func() []mutationType {
					panic("fetchMutationTypes is undefined")
				},
				fetchFilesToMutate: func() []filepath {
					fmt.Println("Fetching files to mutate")
					filepaths := fetchFilesToMutate(&fetchFilesDeps{
						fetchPathsToMutate: func() []filepath { panic("fetchPathsToMutate is undefined") },
						splitFilesAndDirs: func(paths []filepath) (files, dirs []filepath) {
							panic("splitFilesAndDirs is undefined")
						},
						recursivelyExpandDirectories: func(dirs []filepath) (files []filepath) {
							panic("recursivelyExpandDirectories is undefined")
						},
						filterToGoFiles: func(files []filepath) (goFiles []filepath) {
							panic("filterToGoFiles is undefined")
						},
					})
					fmt.Printf("Files to mutate: %v\n", filepaths)

					return filepaths
				},
				testFileMutation: func(filepath, []mutationType) bool {
					panic("testFileMutation is undefined")
				},
			}) {
				fmt.Println("testing mutations failed")
				return false
			}

			fmt.Println("testing mutations passed")

			return true
		},
	}
}

func run(deps *runDeps) bool {
	return deps.pretest() && deps.testMutations()
}

func pretest(deps *pretestDeps) bool {
	c := deps.fetchTestCommand()
	if len(c) == 0 {
		return false
	}

	return deps.runTestCommand(c)
}

func testMutations(deps *testMutationsDeps) bool {
	mutationTypes := deps.fetchMutationTypes()

	if len(mutationTypes) == 0 {
		return false
	}

	filepaths := deps.fetchFilesToMutate()

	if len(filepaths) == 0 {
		return false
	}

	for _, fp := range filepaths {
		if !deps.testFileMutation(fp, mutationTypes) {
			return false
		}
	}

	return true
}

func fetchFilesToMutate(deps *fetchFilesDeps) (filesToMutate []filepath) {
	paths := deps.fetchPathsToMutate()
	files, dirs := deps.splitFilesAndDirs(paths)
	expandedFiles := deps.recursivelyExpandDirectories(dirs)
	allFiles := combine(files, expandedFiles)

	return deps.filterToGoFiles(allFiles)
}

func combine(a, b []filepath) []filepath {
	combined := make([]filepath, len(a)+len(b))
	copy(combined, a)

	for i, item := range b {
		combined[len(a)+i] = item
	}

	return combined
}

type (
	runDeps struct {
		pretest       func() bool
		testMutations func() bool
	}
	pretestDeps struct {
		fetchTestCommand func() command
		runTestCommand   func(command) bool
	}
	command           string
	testMutationsDeps struct {
		fetchMutationTypes func() []mutationType
		fetchFilesToMutate func() []filepath
		testFileMutation   func(filepath, []mutationType) bool
	}
	mutationType struct {
		Name string
	}
	filepath       string
	fetchFilesDeps struct {
		fetchPathsToMutate           func() []filepath
		splitFilesAndDirs            func(paths []filepath) (files, dirs []filepath)
		recursivelyExpandDirectories func(dirs []filepath) (files []filepath)
		filterToGoFiles              func(files []filepath) (goFiles []filepath)
	}
)
