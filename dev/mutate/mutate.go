// Package mutate provides mutation testing functionality.
package main

import "fmt"

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

func main() {
	runner{
		announceMutationTesting: func() { fmt.Printf("Beginning mutation tests\n") },
		testCLICommand: func() bool {
			fmt.Printf("Mock test of CLI Command")
			return true
		},
	}.run()
}

type (
	announceMutationTestingFunc func()
	testCLICommandFunc          func() bool
	runner                      struct {
		announceMutationTesting announceMutationTestingFunc
		testCLICommand          testCLICommandFunc
	}
)

func (r runner) run() {
	r.announceMutationTesting()
	r.testCLICommand()
}
