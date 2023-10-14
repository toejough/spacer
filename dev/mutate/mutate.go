// Package mutate provides mutation testing functionality.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
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
func IOfetchPretestCommand() []string {
	var args struct {
		PretestCommand []string `arg:"positional,required"`
	}

	arg.MustParse(&args)

	return args.PretestCommand
}

func IOrunSubprocess(command []string) bool {
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

func IOexit(code int) {
	os.Exit(code)
}

func getFuncName(f any) string {
	// docs say to use UnsafePointer explicitly instead of Pointer()
	// https://pkg.Pgo.dev/reflect@go1.21.1#Value.Pointer
	return runtime.FuncForPC(uintptr(reflect.ValueOf(f).UnsafePointer())).Name()
}

func debug(function any, args ...any) func(returns ...any) {
	// get function name
	name := getFuncName(function)
	parts := strings.Split(name, ".")
	name = parts[len(parts)-1]
	name = strings.TrimSuffix(name, "-fm")

	if len(args) == 0 {
		fmt.Printf("Called %s...\n", name)
	} else {
		fmt.Printf("Called %s with %v...\n", name, args)
	}
	// print function name & args
	// return a func that prints the function name and returns
	return func(returns ...any) {
		dereferenced := make([]any, len(returns))
		for i := range returns {
			dereferenced[i] = reflect.ValueOf(returns[i]).Elem().Interface()
		}

		if len(returns) == 0 {
			fmt.Printf("...%s completed\n", name)
		} else {
			fmt.Printf("...%s completed with %v\n", name, dereferenced)
		}
	}
}

func unimplemented() string {
	// get function name
	pc := make([]uintptr, 1)
	callsToSkip := 2
	n := runtime.Callers(callsToSkip, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	function := frame.Func
	name := function.Name()

	return fmt.Sprintf("%s is unimplemented\n", name)
}

// Dependency implementations for tested functions.
type prodPretestDeps struct{}

func (pd *prodPretestDeps) fetchPretestCommand() (command []string) {
	defer debug(pd.fetchPretestCommand)(&command)
	return IOfetchPretestCommand()
}

func (pd *prodPretestDeps) runSubprocess(command []string) (result bool) {
	defer debug(pd.runSubprocess)(&result)
	return IOrunSubprocess(command)
}

type prodRunDeps struct{}

func (rd *prodRunDeps) pretest() (result bool) {
	defer debug(rd.pretest)(&result)
	return pretest(&prodPretestDeps{})
}

func (rd *prodRunDeps) testMutations() (result bool) {
	defer debug(rd.testMutations)(&result)
	panic(unimplemented())
}

func (rd *prodRunDeps) exit(code int) {
	defer debug(rd.exit, code)()
	IOexit(code)
}
