// Package mutate provides mutation testing functionality.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
func commmonPrintStarting(fname string) func(string) {
	fmt.Printf("%s is starting...\n", fname)

	return func(result string) {
		fmt.Printf("...%s completed with %s\n", fname, result)
	}
}

// TODO: test the overall flow here - I care about the fetched command being printed to the UI.
func fetchPretestCommand() []string {
	done := commmonPrintStarting("fetchPretestCommand")
	var args struct {
		PretestCommand []string `arg:"positional,required"`
	}

	defer func() { done(strings.Join(args.PretestCommand, " ")) }()
	arg.MustParse(&args)

	return args.PretestCommand
}

// TODO: test the overall flow here - I care about errors being printed to the UI.
func runSubprocess(command []string) bool {
	done := commmonPrintStarting("runSubprocess")
	var cmd string
	var args []string

	if len(command) >= 1 {
		cmd = command[0]
	}
	if len(command) >= 2 {
		args = command[1:]
	}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		exitErr := &exec.ExitError{}
		if errors.As(err, &exitErr) {
			done(string(exitErr.Stderr))
			return false
		}
		done(err.Error())
		return false
	}
	done("Success")
	return true
}

func exit(code int) {
	os.Exit(code)
}

// Dependency implementations for tested functions.
type prodPretestDeps struct{}

func (pd *prodPretestDeps) printStarting(fname string) func(string) {
	return commmonPrintStarting(fname)
}
func (pd *prodPretestDeps) fetchPretestCommand() []string       { return fetchPretestCommand() }
func (pd *prodPretestDeps) runSubprocess(command []string) bool { return runSubprocess(command) }

type prodRunDeps struct{}

func (rd *prodRunDeps) printStarting(fname string) func(string) { return commmonPrintStarting(fname) }
func (rd *prodRunDeps) pretest() bool                           { return pretest(&prodPretestDeps{}) }
func (rd *prodRunDeps) testMutations() bool                     { panic("unimplemented") }
func (rd *prodRunDeps) exit(code int)                           { exit(code) }
