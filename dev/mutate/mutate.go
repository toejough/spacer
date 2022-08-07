// Package mutate provides mutation testing functionality.
package main

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

func main() {
}

type (
	announceMutationTestingFunc func()
	runner                      struct {
		announceMutationTesting announceMutationTestingFunc
	}
)

func (r runner) run() {
	r.announceMutationTesting()
}
