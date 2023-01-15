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

	if run(&runDeps{
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
			panic("testMutations not implemented")
		},
	}) {
		fmt.Println("Mutation testing passed")
		os.Exit(0)
	} else {
		fmt.Println("Mutation testing failed")
		os.Exit(1)
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
    return false
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
	command     string
    testMutationsDeps struct {
        fetchMutationTypes func() []mutationType
        fetchFilesToMutate func() []filepath
        testFileMutation func(filepath, mutationType) bool
    }
    mutationType struct{}
    filepath string
)
