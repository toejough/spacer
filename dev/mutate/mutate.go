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
				fetchTestCommand: func() (command, error) {
					fmt.Println("Fetching test command")
					if len(os.Args) < 2 { //nolint:gomnd
						return "", fmt.Errorf("no test command provided on CLI") //nolint:goerr113
					}
					c := os.Args[1]
					fmt.Printf("Fetched '%s' as the command\n", c)

					return command(c), nil
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
	c, err := deps.fetchTestCommand()
	if err != nil {
		return false
	}

	return deps.runTestCommand(c)
}

type (
	runDeps struct {
		pretest       func() bool
		testMutations func() bool
	}
	command     string
	pretestDeps struct {
		fetchTestCommand func() (command, error)
		runTestCommand   func(command) bool
	}
)
