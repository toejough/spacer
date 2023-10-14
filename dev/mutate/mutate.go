// Package mutate provides mutation testing functionality.
package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/alexflint/go-arg"
)

// Mutate. Based loosely on:
// * https://mutmut.readthedocs.io/en/latest/
// * https://github.com/zimmski/go-mutesting

// Would like to be able to say "replace bool returns with their opposites"
// Would like to cache candidates and results

// main runs the program and exits with 0 on success, 1 on failure, 2 on any kind of runtime failure.
func main() {
	run(&prodRunDeps{})
}

// Untested io funcs which need integration testing rather than unit testing.
func fetchPretestCommand() []string {
	var args struct {
		PretestCommand []string `arg:"positional,required"`
	}

	arg.MustParse(&args)

	return args.PretestCommand
}

func runSubprocess(command []string) bool {
	var (
		cmd  string
		args []string
	)

	// len has to be over 1 or there's no command
	if len(command) >= 1 {
		cmd = command[0]
	}

	// len has to be over 2 or there're no args
	if len(command) >= 2 { //nolint:gomnd
		args = command[1:]
	}

	if err := exec.Command(cmd, args...).Run(); err != nil {
		exitErr := new(exec.ExitError)
		if errors.As(err, &exitErr) {
			return false
		}

		return false
	}

	return true
}

func exit(code int) {
	os.Exit(code)
}

// TODO: write a debug function that:
// prints the name of the function with args
// returns a func to be called when the function is done
// that func prints the name of the function and the return values

// Dependency implementations for tested functions.
type prodPretestDeps struct{}

// TODO: make the UI stuff happen here, actually. announcing stuff is starting/done with what result.
func (pd *prodPretestDeps) fetchPretestCommand() []string {
	return fetchPretestCommand()
}
func (pd *prodPretestDeps) runSubprocess(command []string) bool { return runSubprocess(command) }

type prodRunDeps struct{}

func (rd *prodRunDeps) pretest() bool       { return pretest(&prodPretestDeps{}) }
func (rd *prodRunDeps) testMutations() bool { panic("unimplemented") }
func (rd *prodRunDeps) exit(code int)       { exit(code) }
