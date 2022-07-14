package main

import (
	"fmt"
	"spacer/dev/protest"
	"testing"
)

type diffBool bool

func (e diffBool) Diff(b diffBool) string {
	if bool(e) != bool(b) {
		return fmt.Sprintf("%t != %t", bool(e), bool(b))
	}

	return ""
}

func TestRun(t *testing.T) {
	t.Parallel()

	for _, mutationReturn := range []bool{true, false} { //nolint:paralleltest
		// linter can't tell we're using the range value in the test
		// Given a return from the mutation func
		mReturn := mutationReturn

		t.Run(fmt.Sprintf("%t", mutationReturn), func(t *testing.T) {
			t.Parallel()

			actualCalls := protest.Stack[string]{}
			var reportArgs diffBool
			var exitArgs diffBool

			// Given dependencies
			mutate := func() bool {
				actualCalls.Push("mutate")

				return mReturn
			}
			report := func(r bool) {
				actualCalls.Push("report")
				reportArgs = diffBool(r)
			}
			exit := func(r bool) {
				actualCalls.Push("exit")
				exitArgs = diffBool(r)
			}

			// When run is called
			run(mutate, report, exit)

			// Then the functions are called in the right order with the right args
			protest.RequireCall(t, "mutate", actualCalls.MustPop(t))

			protest.RequireCall(t, "report", actualCalls.MustPop(t))
			protest.RequireArgs(t, diffBool(mReturn), reportArgs)

			protest.RequireCall(t, "exit", actualCalls.MustPop(t))
			protest.RequireArgs(t, diffBool(mReturn), exitArgs)

			protest.RequireEmpty(t, actualCalls)
		})
	}
}
