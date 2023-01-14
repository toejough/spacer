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
	if run(&runDeps{
		announceStarting: func() { fmt.Println("Starting mutation testing") },
		pretest: func() bool {
			return pretest(&pretestDeps{
				announcePretest: func() { fmt.Println("Starting pretesting") },
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
				announcePretestResults: func(b bool) { fmt.Printf("Pretest passed? %t\n", b) },
			})
		},
		testMutations: func() bool {
			panic("testMutations not implemented")
		},
		announceEnding: func(b bool) { fmt.Printf("Mutation testing passed? %t\n", b) },
	}) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func run(deps *runDeps) bool {
	deps.announceStarting()
	passes := deps.pretest() && deps.testMutations()
	deps.announceEnding(passes)

	return passes
}

func pretest(deps *pretestDeps) bool {
	deps.announcePretest()

	c, err := deps.fetchTestCommand()
	if err != nil {
		return false
	}

	result := deps.runTestCommand(c)
	deps.announcePretestResults(result)

	return result
}

type (
	runDeps struct {
		announceStarting func()
		pretest          func() bool
		testMutations    func() bool
		announceEnding   func(bool)
	}
	command     string
	pretestDeps struct {
		announcePretest        func()
		fetchTestCommand       func() (command, error)
		runTestCommand         func(command) bool
		announcePretestResults func(bool)
	}
)
